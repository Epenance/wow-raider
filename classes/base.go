package classes

import (
	"fmt"
	"github.com/kbinani/screenshot"
	"github.com/moutend/go-hook/pkg/keyboard"
	"github.com/moutend/go-hook/pkg/types"
	"image"
	"image/png"
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
}

func (c *BaseClass) Uninit() {
	keyboard.Uninstall()
}

func (c *BaseClass) Init() error {
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
	// defer keyboard.Uninstall()
	// defer fmt.Println("Listening for keyboard events")

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

func (c *BaseClass) CastSpell(spell string, customDelay ...time.Duration) bool {
	fmt.Println("Casting Spell: " + spell)
	return true
}

func PrintFields(val reflect.Value) {
	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)

		// Check if the field is a nested struct
		if valueField.Kind() == reflect.Struct {
			PrintFields(valueField)
		} else {
			fmt.Printf("%s: %v\n", typeField.Name, valueField.Interface())
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

func (c *BaseClass) SetWeakAuraSize() {
	var colorToFind = util.COLORS["BLUE"]
	bounds := c.GameScreenshot.Bounds()
	var startX, startY, endX, endY int = -1, -1, -1, -1

	// Find the first BLUE pixel
	for y := bounds.Min.Y; y <= bounds.Max.Y && startX == -1; y++ {
		for x := bounds.Min.X; x <= bounds.Max.X; x++ {
			if util.IsColor(colorToFind, c.GameScreenshot, x, y) {
				startX, startY = x, y
				break
			}
		}
	}

	if startX == -1 { // No BLUE pixel found
		return
	}

	// Find the last BLUE pixel to the right from startX, startY
	endX = startX
	for x := startX; x <= bounds.Max.X; x++ {
		if util.IsColor(colorToFind, c.GameScreenshot, x, startY) {
			endX = x
		} else {
			break
		}
	}

	// Find the last BLUE pixel downward from startX, endY
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

}
