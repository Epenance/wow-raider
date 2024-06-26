package classes

import (
	"fmt"
	"github.com/epenance/virtual_keyboard"
	"github.com/gdamore/tcell/v2"
	"github.com/kbinani/screenshot"
	"github.com/moutend/go-hook/pkg/keyboard"
	"github.com/moutend/go-hook/pkg/types"
	"github.com/rivo/tview"
	"gopkg.in/yaml.v3"
	"image"
	"image/png"
	"io/ioutil"
	"os"
	"os/signal"
	"reflect"
	"sort"
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
	IsCasting        bool
	ChatOpen         bool
	OnGlobalCooldown bool
}

type WeakAura struct {
	width  int
	height int
}

type BaseClass struct {
	HWND             window.HWND
	TView            *tview.Application
	TViewTables      map[string]*tview.Table
	TViewTableValues map[string]map[string]TableCellValue
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

type ConfigKeybinding struct {
	Key      string `yaml:"key"`
	HasCtrl  bool   `yaml:"ctrl,omitempty"`
	HasShift bool   `yaml:"shift,omitempty"`
	HasAlt   bool   `yaml:"alt,omitempty"`
}

type Config struct {
	Keys map[string]ConfigKeybinding `yaml:"keys"`
}

type TableCellValue struct {
	ZIndex     int
	NameColor  tcell.Color
	Value      string
	ValueColor tcell.Color
}

type KeyListener struct {
	Key      types.VKCode
	Function func()
}

func (c *BaseClass) Uninit() {
	keyboard.Uninstall()
}

func (c *BaseClass) Init(listeners []KeyListener) error {
	c.TViewTables = make(map[string]*tview.Table)

	c.TViewTableValues = make(map[string]map[string]TableCellValue)
	c.TViewTableValues["options"] = make(map[string]TableCellValue)
	c.TViewTableValues["state"] = make(map[string]TableCellValue)

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

				// Add listeners
				for _, listener := range listeners {
					if event.VKCode == listener.Key && event.Message == types.WM_KEYDOWN {
						listener.Function()
					}
				}
			case <-interrupt:
				// If an interrupt signal is received, stop the program.
				util.Log("Terminating Program...")
				c.InterruptProgram = true
				return
			}
		}
	}()

	c.TView = tview.NewApplication()

	options := tview.NewTable().SetSelectedStyle(tcell.Style{}.Foreground(tcell.ColorBlack).Background(tcell.ColorWhite))
	options.SetTitle("Options").SetBorder(true)

	stateTable := tview.NewTable().SetSelectedStyle(tcell.Style{}.Foreground(tcell.ColorBlack).Background(tcell.ColorWhite))
	stateTable.SetTitle("State").SetBorder(true)

	c.TViewTables["options"] = options
	c.TViewTables["state"] = stateTable

	logs := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWordWrap(true).
		SetScrollable(true)

	logs.SetChangedFunc(func() {
		c.TView.Draw()
		logs.ScrollToEnd() // Automatically scroll to the end
	})

	util.SetWriter(logs)

	go func() {
		flex := tview.NewFlex().
			AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
				AddItem(options, 0, 1, false).
				AddItem(stateTable, 0, 2, false), 0, 1, false).
			AddItem(logs, 0, 2, false)
		if err := c.TView.SetRoot(flex, true).EnableMouse(true).Run(); err != nil {
			panic(err)
		}

		// Closing down the app when the program is interrupted (ctrl + c)
		interrupt <- syscall.SIGINT
	}()

	util.Log("Initialized " + c.Spec + " " + c.Class)

	c.RunProgram = true

	return nil
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

func (c *BaseClass) CastSpell(spell string, customDelay ...time.Duration) {
	var delay time.Duration
	if len(customDelay) > 0 {
		delay = customDelay[0]
	} else {
		delay = 150
	}

	kb, err := virtual_keyboard.NewKeyBonding()
	kb.AddHWND(uintptr(c.HWND))
	if err != nil {
		panic(err)
	}

	// Check if the spell exists in the keybindings
	if _, ok := c.Keybindings[spell]; !ok {
		util.Log("Keybinding not found for: " + spell)
		return
	}

	kb.SetKeys(c.Keybindings[spell].Key)

	if c.Keybindings[spell].HasShift {
		kb.HasSHIFT(true)
	}

	if c.Keybindings[spell].HasCtrl {
		kb.HasCTRL(true)
	}

	if c.Keybindings[spell].HasAlt {
		kb.HasALT(true)
	}
	util.Log("Casting: " + spell)
	err = kb.Launch()
	if err != nil {
		panic(err)
	}

	time.Sleep(delay * time.Millisecond)
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
	c.State.IsAlive = c.CheckColor(util.PURPLE, 0, 5)
	c.State.InCombat = c.CheckColor(util.PURPLE, 10, 0)
	c.State.IsMounted = c.CheckColor(util.PURPLE, 5, 5)
	c.State.ChatOpen = c.CheckColor(util.PURPLE, 15, 5)
	c.State.OnGlobalCooldown = c.CheckColor(util.PURPLE, 40, 5)
	c.State.IsCasting = c.CheckColor(util.PURPLE, 10, 5)
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

func (c *BaseClass) UpdateTables() {
	c.TView.QueueUpdateDraw(func() {
		options := c.TViewTables["options"]
		stateTable := c.TViewTables["state"]

		optionValues := c.TViewTableValues["options"]
		stateValues := c.TViewTableValues["state"]

		optionValues["Run Program"] = TableCellValue{ZIndex: 10, NameColor: tcell.ColorWhite, Value: fmt.Sprintf("%t", c.RunProgram), ValueColor: util.GetColor(c.RunProgram, tcell.ColorGreen, tcell.ColorRed)}
		optionValues["Use cooldowns"] = TableCellValue{ZIndex: 20, NameColor: tcell.ColorWhite, Value: fmt.Sprintf("%t", c.PopCooldowns), ValueColor: util.GetColor(c.PopCooldowns, tcell.ColorGreen, tcell.ColorRed)}
		optionValues["Rotate cooldowns"] = TableCellValue{ZIndex: 30, NameColor: tcell.ColorWhite, Value: fmt.Sprintf("%t", c.ForceCooldowns), ValueColor: util.GetColor(c.ForceCooldowns, tcell.ColorGreen, tcell.ColorRed)}

		stateValues["Is Alive"] = TableCellValue{ZIndex: 999, NameColor: tcell.ColorWhite, Value: fmt.Sprintf("%t", c.State.IsAlive), ValueColor: util.GetColor(c.State.IsAlive, tcell.ColorGreen, tcell.ColorRed)}
		stateValues["In Combat"] = TableCellValue{ZIndex: 998, NameColor: tcell.ColorWhite, Value: fmt.Sprintf("%t", c.State.InCombat), ValueColor: util.GetColor(c.State.InCombat, tcell.ColorGreen, tcell.ColorRed)}
		stateValues["Is Mounted"] = TableCellValue{ZIndex: 997, NameColor: tcell.ColorWhite, Value: fmt.Sprintf("%t", c.State.IsMounted), ValueColor: util.GetColor(c.State.IsMounted, tcell.ColorGreen, tcell.ColorRed)}
		stateValues["Chat Open"] = TableCellValue{ZIndex: 996, NameColor: tcell.ColorWhite, Value: fmt.Sprintf("%t", c.State.ChatOpen), ValueColor: util.GetColor(c.State.ChatOpen, tcell.ColorGreen, tcell.ColorRed)}
		stateValues["On Global Cooldown"] = TableCellValue{ZIndex: 100, NameColor: tcell.ColorWhite, Value: fmt.Sprintf("%t", c.State.OnGlobalCooldown), ValueColor: util.GetColor(c.State.OnGlobalCooldown, tcell.ColorGreen, tcell.ColorRed)}

		// Remap the option values to a sortable slice
		sortedOptionValues := make([]struct {
			Key   string
			Value TableCellValue
		}, 0, len(optionValues))

		for k, v := range optionValues {
			sortedOptionValues = append(sortedOptionValues, struct {
				Key   string
				Value TableCellValue
			}{Key: k, Value: v})
		}

		// Sort slice by ZIndex, then alphabetically by Key
		sort.Slice(sortedOptionValues, func(i, j int) bool {
			if sortedOptionValues[i].Value.ZIndex == sortedOptionValues[j].Value.ZIndex {
				return sortedOptionValues[i].Key < sortedOptionValues[j].Key // Alphabetical order if ZIndex are equal
			}
			return sortedOptionValues[i].Value.ZIndex < sortedOptionValues[j].Value.ZIndex
		})

		for i, kv := range sortedOptionValues {
			options.SetCell(i, 0, tview.NewTableCell(kv.Key).SetTextColor(kv.Value.NameColor).SetAlign(tview.AlignLeft))
			options.SetCell(i, 1, tview.NewTableCell(kv.Value.Value).SetTextColor(kv.Value.ValueColor).SetAlign(tview.AlignLeft).SetExpansion(1))
		}

		// Remap the state values to a sortable slice
		sortedStateValues := make([]struct {
			Key   string
			Value TableCellValue
		}, 0, len(stateValues))

		for k, v := range stateValues {
			sortedStateValues = append(sortedStateValues, struct {
				Key   string
				Value TableCellValue
			}{Key: k, Value: v})
		}

		// Sort slice by ZIndex, then alphabetically by Key
		sort.Slice(sortedStateValues, func(i, j int) bool {
			if sortedStateValues[i].Value.ZIndex == sortedStateValues[j].Value.ZIndex {
				return sortedStateValues[i].Key < sortedStateValues[j].Key // Alphabetical order if ZIndex are equal
			}
			return sortedStateValues[i].Value.ZIndex < sortedStateValues[j].Value.ZIndex
		})

		for i, kv := range sortedStateValues {
			stateTable.SetCell(i, 0, tview.NewTableCell(kv.Key).SetTextColor(kv.Value.NameColor).SetAlign(tview.AlignLeft))
			stateTable.SetCell(i, 1, tview.NewTableCell(kv.Value.Value).SetTextColor(kv.Value.ValueColor).SetAlign(tview.AlignLeft).SetExpansion(1))
		}
	})
}
