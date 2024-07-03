package shaman

import (
	"fmt"
	"github.com/moutend/go-hook/pkg/types"
	"time"
	"wow-raider/classes"
	"wow-raider/util"
)

type EnhancementState struct {
	ShamanState
	LavaLashAvailable     bool
	StormstrikeAvailable  bool
	MaelstromWeaponStacks int
}

type Enhancement struct {
	Shaman
	State EnhancementState
}

func (c *Enhancement) Run() {

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

func (c *Enhancement) Init() error {
	c.Spec = "Enhancement"

	// []classes.KeyListener
	listeners := []classes.KeyListener{}

	listeners = append(listeners, classes.KeyListener{Key: types.VK_F4, Function: func() {
		fmt.Println("test")
	}})

	if err := c.Shaman.Init(listeners); err != nil {
		return err
	}

	return nil
}

func (c *Enhancement) SetState() {
	c.Shaman.SetState()

	// Hack because there is no inheritance in Go
	c.SyncState(&c.Shaman.State, &c.State)

	c.State.LavaLashAvailable = c.CheckColor(util.BLUE, 20, 0)
	c.State.StormstrikeAvailable = c.CheckColor(util.BLUE, 25, 0)

	maelstromWeaponsFive := c.CheckColor(util.GREEN, 35, 0)

	if maelstromWeaponsFive {
		c.State.MaelstromWeaponStacks = 5
	} else {
		c.State.MaelstromWeaponStacks = 0
	}
}

func (c *Enhancement) Rotation() {

	if c.State.ChatOpen {
		return
	}

	state := c.State
	combatAliveAndNotMounted := state.IsAlive && !state.IsMounted && state.InCombat

	if state.IsAlive && !state.IsMounted && !state.OnGlobalCooldown && state.WindfuryMissing {
		c.CastSpell("Windfury Weapon")
		return
	}

	if state.IsAlive && !state.IsMounted && !state.OnGlobalCooldown && state.FlametongueMissing && !state.WindfuryMissing {
		c.CastSpell("Flametongue Weapon")
		return
	}

	if combatAliveAndNotMounted && !state.OnGlobalCooldown && state.FlameShockAvailable && !state.FlameShockDotActive {
		c.CastSpell("Flame Shock")
		return
	}

	if combatAliveAndNotMounted && !state.OnGlobalCooldown && state.FireNovaAvailable && state.FlameshockDotsActive >= 3 {
		c.CastSpell("Fire Nova")
		return
	}

	if combatAliveAndNotMounted && !state.OnGlobalCooldown && state.StormstrikeAvailable {
		c.CastSpell("Stormstrike")
		return
	}

	if combatAliveAndNotMounted && !state.OnGlobalCooldown && state.LavaLashAvailable {
		c.CastSpell("Lava Lash")
		return
	}

	if combatAliveAndNotMounted && !state.OnGlobalCooldown && state.MaelstromWeaponStacks == 5 {
		c.CastSpell("Lightning Bolt", 500)
		return
	}

	if combatAliveAndNotMounted && !state.OnGlobalCooldown && state.EarthShockAvailable && state.FlameShockDotActive {
		c.CastSpell("Earth Shock")
		return
	}

	if !state.IsMounted && !state.OnGlobalCooldown && state.LightningShieldMissing {
		c.CastSpell("Lightning Shield")
		return
	}
}

func (c *Enhancement) UpdateTables() {
	// optionValues := c.TViewTableValues["options"]
	// stateValues := c.TViewTableValues["state"]

	c.Shaman.UpdateTables()
}
