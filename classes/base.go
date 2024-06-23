package classes

import (
	"fmt"
	"github.com/moutend/go-hook/pkg/keyboard"
	"github.com/moutend/go-hook/pkg/types"
	"os"
	"os/signal"
	"reflect"
	"syscall"
	"time"
	"wow-raider/util"
)

type BaseState struct {
	IsLoading        bool
	IsAlive          bool
	InCombat         bool
	IsMounted        bool
	ChatOpen         bool
	OnGlobalCooldown bool
}

type BaseClass struct {
	HWND             uintptr
	RunProgram       bool
	PopCooldowns     bool
	ForceCooldowns   bool
	InterruptProgram bool
	State            BaseState
}

func (c *BaseClass) Init() error {
	// Setup interrupt handling
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	// Create a channel to receive keyboard events.
	events := make(chan types.KeyboardEvent, 1024)

	if err := keyboard.Install(nil, events); err != nil {
		panic(err)
	}
	defer keyboard.Uninstall()

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
			case <-interrupt:
				// If an interrupt signal is received, stop the program.
				util.Log("Terminating Program...")
				c.InterruptProgram = true
				return
			}
		}
	}()

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
