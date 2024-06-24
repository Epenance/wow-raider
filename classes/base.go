package classes

import (
	"fmt"
	"github.com/kbinani/screenshot"
	"github.com/moutend/go-hook/pkg/keyboard"
	"github.com/moutend/go-hook/pkg/types"
	"gopkg.in/yaml.v3"
	"image"
	"image/png"
	"io/ioutil"
	"os"
	"os/signal"
	"reflect"
	"syscall"
	"time"
	"wow-raider/util"
	"wow-raider/window"
)

type BaseState struct {
	IsLoading        bool
	IsAlive          bool
	InCombat         bool
	IsMounted        bool
	ChatOpen         bool
	OnGlobalCooldown bool
}

type WeakAura struct {
	width  int
	height int
}

type BaseClass struct {
	HWND             window.HWND
	Class            string
	Spec             string
	RunProgram       bool
	PopCooldowns     bool
	ForceCooldowns   bool
	InterruptProgram bool
	State            BaseState
	GameScreenshot   image.Image
	WeakAura         WeakAura
	Keybindings      map[string]ConfigKeybinding
}

func (c *BaseClass) Uninit() {
	keyboard.Uninstall()
}

func (c *BaseClass) Init() error {
	err := c.LoadConfig()
	if err != nil {
		return fmt.Errorf("config not found")
	}

	hwnd := window.FindWindowByTitle("World of Warcraft")

	if hwnd == 0 {
		return fmt.Errorf("window not found")
	}

	c.HWND = hwnd

	// Setup interrupt handling
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	// Create a channel to receive keyboard events.
	events := make(chan types.KeyboardEvent, 1024)

	if err := keyboard.Install(nil, events); err != nil {
		panic(err)
	}

	go func() {

		for {
			select {
			case event := <-events:
				if event.VKCode == types.VK_PAUSE && event.Message == types.WM_KEYDOWN {
					if c.RunProgram {
						util.Log("Program paused")
					} else {
						util.Log("Program resumed")
					}
					c.RunProgram = !c.RunProgram
				}

				if event.VKCode == types.VK_END && event.Message == types.WM_KEYDOWN {
					if c.PopCooldowns {
						util.Log("Cancelling CD's")
					} else {
						util.Log("Popping CD's on next available opportunity")
					}
					c.PopCooldowns = !c.PopCooldowns
				}

				if event.VKCode == types.VK_F1 && event.Message == types.WM_KEYDOWN {
					if c.ForceCooldowns {
						util.Log("Stopping CD forcing")
					} else {
						util.Log("Casting CD's on cooldown until stopped")
					}
					c.ForceCooldowns = !c.ForceCooldowns
				}

				if event.VKCode == types.VK_F2 && event.Message == types.WM_KEYDOWN {
					util.Log("Saving screenshot")
					c.SaveScreenshot()
				}
			case <-interrupt:
				// If an interrupt signal is received, stop the program.
				util.Log("Terminating Program...")
				c.InterruptProgram = true
				return
			}
		}
	}()

	util.Log("Initialized " + c.Spec + " " + c.Class)

	c.RunProgram = true

	return nil
}

type ConfigKeybinding struct {
	Key      string `yaml:"key"`
	HasCtrl  bool   `yaml:"ctrl,omitempty"`
	HasShift bool   `yaml:"shift,omitempty"`
	HasAlt   bool   `yaml:"alt,omitempty"`
}

type Config struct {
	Keys map[string]ConfigKeybinding `yaml:"keys"`
}

func (c *BaseClass) LoadConfig() error {
	// Read the YAML file
	data, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		return err
	}

	// Decode the YAML data into Config
	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return err
	}

	// Assign the decoded keys to the BaseClass Keybindings
	c.Keybindings = cfg.Keys

	return nil

}

func (c *BaseClass) CastSpell(spell string, customDelay ...time.Duration) bool {
	fmt.Println("Casting Spell: " + spell)
	return true
}

func PrintFields(val reflect.Value) {
	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)

		// Check if the field is exported
		if valueField.CanInterface() {
			// Check if the field is a nested struct
			if valueField.Kind() == reflect.Struct {
				PrintFields(valueField)
			} else {
				fmt.Printf("%s: %v\n", typeField.Name, valueField.Interface())
			}
		} else {
			// Ignore the field
		}
	}
}

func (c *BaseClass) PrintState() {
	// Use reflect.ValueOf(c).Elem() to get the correct value
	PrintFields(reflect.ValueOf(c.State))
}

func (c *BaseClass) CaptureGame() error {
	if c.HWND == 0 {
		return fmt.Errorf("hwnd not found")
	}

	dimensions, err := window.DwmGetWindowBounds(c.HWND)
	if err != nil {
		return fmt.Errorf("window not found")
	}

	b := image.Rect(
		int(dimensions.Left+1), int(dimensions.Top),
		int(dimensions.Right-1), int(dimensions.Bottom-1))

	img, err := screenshot.CaptureRect(b)

	if err != nil {
		return err
	}

	c.GameScreenshot = img

	c.SetWeakAuraSize()

	return nil
}

func (c *BaseClass) CheckColor(color util.RGB, x, y int) bool {
	x, y = c.GetActualPosition(x, y)
	return util.IsColor(color, c.GameScreenshot, x, y)
}

func (c *BaseClass) GetColor(x, y int) util.RGB {
	x, y = c.GetActualPosition(x, y)
	r, g, b := util.PixelColorAt(c.GameScreenshot, x, y)
	return util.RGB{R: r, G: g, B: b}
}

func (c *BaseClass) GetActualPosition(x, y int) (int, int) {
	inGameX, inGameY := 5, 5
	waX, waY := c.WeakAura.width, c.WeakAura.height

	// Fallback Sizes
	if waX == 0 {
		waX = 5
	}

	if waY == 0 {
		waY = 5
	}

	// Get the multiplier for X & Y size
	multiplierX := x / inGameX
	multiplierY := y / inGameY

	// Real Pixel Positions
	x = waX * multiplierX
	y = waY * multiplierY

	// Now add the center of the Weak Aura sizes
	x += waX / 2
	y += waY / 2

	bounds := c.GameScreenshot.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	return width - x, height - y
}

func (c *BaseClass) SetWeakAuraSize() {
	var colorToFind = util.COLORS["BLUE"]
	bounds := c.GameScreenshot.Bounds()
	var startX, startY, endX, endY int = -1, -1, -1, -1

	// Find the first matching pixel
	for y := bounds.Min.Y; y <= bounds.Max.Y && startX == -1; y++ {
		for x := bounds.Min.X; x <= bounds.Max.X; x++ {
			if util.IsColor(colorToFind, c.GameScreenshot, x, y) {
				startX, startY = x, y
				break
			}
		}
	}

	if startX == -1 { // No matching pixel found
		return
	}

	// Find the last matching pixel to the right from startX, startY
	endX = startX
	for x := startX; x <= bounds.Max.X; x++ {
		if util.IsColor(colorToFind, c.GameScreenshot, x, startY) {
			endX = x
		} else {
			break
		}
	}

	// Find the last matching pixel downward from startX, endY
	endY = startY
	for y := startY; y <= bounds.Max.Y; y++ {
		if util.IsColor(colorToFind, c.GameScreenshot, startX, y) {
			endY = y
		} else {
			break
		}
	}

	// Calculate dimensions
	if endX != -1 && endY != -1 {
		c.WeakAura.width = endX - startX + 1
		c.WeakAura.height = endY - startY + 1
	}
}

func (c *BaseClass) SaveScreenshot() {
	fileName := "hello.png"
	file, _ := os.Create(fileName)
	err := png.Encode(file, c.GameScreenshot)
	if err != nil {
		return
	}
	err = file.Close()
	if err != nil {
		return
	}
}

func (c *BaseClass) SetState() {
	c.State.IsAlive = !c.CheckColor(util.PURPLE, 0, 5)
	c.State.InCombat = c.CheckColor(util.PURPLE, 10, 0)
	c.State.IsMounted = c.CheckColor(util.PURPLE, 5, 5)
	c.State.ChatOpen = c.CheckColor(util.PURPLE, 15, 5)
	c.State.OnGlobalCooldown = c.CheckColor(util.GREEN, 40, 5)
}

func (c *BaseClass) SyncState(base, target interface{}) {
	baseVal := reflect.ValueOf(base).Elem()
	targetVal := reflect.ValueOf(target).Elem()

	for i := 0; i < baseVal.NumField(); i++ {
		baseField := baseVal.Field(i)
		targetField := targetVal.FieldByName(baseVal.Type().Field(i).Name)

		if targetField.IsValid() && targetField.CanSet() {
			targetField.Set(baseField)
		}
	}
}
