package shaman

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/moutend/go-hook/pkg/types"
	"time"
	"wow-raider/classes"
	"wow-raider/util"
)

type ElementalState struct {
	ShamanState
	IsMoving              bool
	LightningShieldCount  int
	ThunderstormAvailable bool
	LavaBurstAvailable    bool
	LightingBoltAvailable bool
	FlametongueMissing    bool
}

type Elemental struct {
	Shaman
	State            ElementalState
	AoEEnabled       bool
	UseFireElemental bool
	UseBloodlust     bool
}

func (c *Elemental) Run() {

	go func() {
		for !c.InterruptProgram {
			c.UpdateTables()
			time.Sleep(100 * time.Millisecond)
		}
	}()

	frequency := 30                                 // Updates per second
	delay := time.Second / time.Duration(frequency) // Delay between each iteration

	for !c.InterruptProgram {
		startTime := time.Now()

		if c.RunProgram {
			err := c.CaptureGame()

			if err != nil {
				return
			}

			c.SetState()

			c.Rotation()
		}

		elapsed := time.Since(startTime) // Time spend in this iteration
		if elapsed < delay {
			time.Sleep(delay - elapsed) // Delay for the remaining time
		}
	}

}

func (c *Elemental) Init() error {
	c.Spec = "Elemental"

	// []classes.KeyListener
	listeners := []classes.KeyListener{}

	listeners = append(listeners, classes.KeyListener{Key: types.VK_F2, Function: func() {
		c.AoEEnabled = !c.AoEEnabled
	}})

	listeners = append(listeners, classes.KeyListener{Key: types.VK_F3, Function: func() {
		c.UseFireElemental = !c.UseFireElemental
	}})

	listeners = append(listeners, classes.KeyListener{Key: types.VK_F4, Function: func() {
		c.UseBloodlust = !c.UseBloodlust
	}})

	if err := c.Shaman.Init(listeners); err != nil {
		return err
	}

	return nil
}

func (c *Elemental) SetState() {
	c.Shaman.SetState()

	// Hack because there is no inheritance in Go
	c.SyncState(&c.Shaman.State, &c.State)

	c.State.LavaBurstAvailable = c.CheckColor(util.BLUE, 35, 0)
	c.State.LightingBoltAvailable = c.CheckColor(util.BLUE, 25, 0)
	c.State.ThunderstormAvailable = c.CheckColor(util.BLUE, 20, 0)
	c.State.IsMoving = c.CheckColor(util.RED, 45, 0)
	c.State.FlametongueMissing = c.CheckColor(util.RED, 25, 5)

	lightningShieldStacks := c.CheckColor(util.GREEN, 40, 0)

	if lightningShieldStacks {
		c.State.LightningShieldCount = 7
	} else {
		c.State.LightningShieldCount = 0
	}
}

func (c *Elemental) Rotation() {

	// TODO:
	// - Add support for Totems
	// - Add support for Cooldowns

	if c.State.ChatOpen {
		return
	}

	state := c.State
	combatAliveAndNotMounted := state.IsAlive && !state.IsMounted && state.InCombat

	if state.IsAlive && !state.IsMounted && !state.OnGlobalCooldown && state.FlametongueMissing && !state.IsCasting {
		c.CastSpell("Flametongue Weapon")
		return
	}

	if state.IsAlive && !state.IsMounted && !state.OnGlobalCooldown && state.LightningShieldMissing {
		c.CastSpell("Lightning Shield")
		return
	}

	if c.AoEEnabled {
		if combatAliveAndNotMounted && !state.OnGlobalCooldown && !state.IsCasting && state.ThunderstormAvailable {
			c.CastSpell("Thunderstorm")
			return
		}

		if combatAliveAndNotMounted && !state.OnGlobalCooldown && !state.IsCasting && state.EarthShockAvailable && state.LightningShieldCount > 6 {
			c.CastSpell("Earth Shock")
			return
		}

		if combatAliveAndNotMounted && !state.OnGlobalCooldown && !state.IsMoving && !state.IsCasting {
			c.CastSpell("Chain Lightning")
			return
		}
	}

	if combatAliveAndNotMounted && !state.OnGlobalCooldown && !state.IsCasting && state.FlameShockAvailable && !state.FlameShockDotActive {
		c.CastSpell("Flame Shock")
		return
	}

	if combatAliveAndNotMounted && !state.OnGlobalCooldown && !state.IsCasting && state.LavaBurstAvailable && !state.IsMoving && (state.FlameShockAvailable || state.FlameShockDotActive) {
		c.CastSpell("Lava Burst")
		return
	}

	if combatAliveAndNotMounted && !state.OnGlobalCooldown && !state.IsCasting && state.EarthShockAvailable && state.LightningShieldCount > 6 {
		c.CastSpell("Earth Shock")
		return
	}

	if combatAliveAndNotMounted && !state.OnGlobalCooldown && !state.IsCasting && state.ThunderstormAvailable {
		c.CastSpell("Thunderstorm")
		return
	}

	if combatAliveAndNotMounted && !state.OnGlobalCooldown && !state.IsCasting {
		c.CastSpell("Lightning Bolt")
		return
	}

}

func (c *Elemental) UpdateTables() {
	optionValues := c.TViewTableValues["options"]
	stateValues := c.TViewTableValues["state"]

	optionValues["AoE Enabled (F2)"] = classes.TableCellValue{ZIndex: 1, NameColor: tcell.ColorWhite, Value: fmt.Sprintf("%t", c.AoEEnabled), ValueColor: util.GetColor(c.AoEEnabled, tcell.ColorGreen, tcell.ColorRed)}
	optionValues["Fire elemental (F3)"] = classes.TableCellValue{ZIndex: 2, NameColor: tcell.ColorWhite, Value: fmt.Sprintf("%t", c.UseFireElemental), ValueColor: util.GetColor(c.UseFireElemental, tcell.ColorGreen, tcell.ColorRed)}
	optionValues["Bloodlust (F4)"] = classes.TableCellValue{ZIndex: 3, NameColor: tcell.ColorWhite, Value: fmt.Sprintf("%t", c.UseBloodlust), ValueColor: util.GetColor(c.UseBloodlust, tcell.ColorGreen, tcell.ColorRed)}
	stateValues["Flametongue Missing"] = classes.TableCellValue{ZIndex: 1, NameColor: tcell.ColorWhite, Value: fmt.Sprintf("%t", c.State.FlametongueMissing), ValueColor: util.GetColor(c.State.FlametongueMissing, tcell.ColorGreen, tcell.ColorRed)}

	c.Shaman.UpdateTables()
}
