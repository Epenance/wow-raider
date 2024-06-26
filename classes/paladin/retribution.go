package paladin

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/moutend/go-hook/pkg/types"
	"time"
	"wow-raider/classes"
	"wow-raider/util"
)

type RetributionState struct {
	PaladinState
	DivineStormAvailable      bool
	ArtOfWarActive            bool
	TemplarsVerdictReady      bool
	DivinePurposeActive       bool
	InquisitionActive         bool
	ZealotryActive            bool
	ZealotryAvailable         bool
	JudgementsOfThePureActive bool
}

type Retribution struct {
	Paladin
	State RetributionState
}

func (c *Retribution) Run() {

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

func (c *Retribution) Init() error {
	c.Spec = "Retribution"

	// []classes.KeyListener
	listeners := []classes.KeyListener{}

	listeners = append(listeners, classes.KeyListener{Key: types.VK_F4, Function: func() {
		fmt.Println("test")
	}})

	if err := c.Paladin.Init(listeners); err != nil {
		return err
	}

	return nil
}

func (c *Retribution) SetState() {
	c.Paladin.SetState()

	// Hack because there is no inheritance in Go
	c.SyncState(&c.Paladin.State, &c.State)

	c.State.DivineStormAvailable = c.CheckColor(util.BLUE, 30, 0)
	c.State.ArtOfWarActive = c.CheckColor(util.BLUE, 20, 5)
	c.State.TemplarsVerdictReady = c.CheckColor(util.BLUE, 50, 5)
	c.State.DivinePurposeActive = c.CheckColor(util.BLUE, 55, 5)
	c.State.InquisitionActive = c.CheckColor(util.GREEN, 55, 0)
	c.State.ZealotryActive = c.CheckColor(util.GREEN, 60, 0)
	c.State.ZealotryAvailable = c.CheckColor(util.BLUE, 65, 0)
	c.State.JudgementsOfThePureActive = c.CheckColor(util.GREEN, 70, 0)
}

func (c *Retribution) Rotation() {
	// Do stuff
	now := time.Now()
	fiveSeconds := 5 * time.Second

	if c.State.ChatOpen {
		return
	}

	state := c.State

	if now.Sub(c.State.LastZealCast) > fiveSeconds {
		// fmt.Println("5 seconds gone by?")
		if state.IsAlive &&
			!state.IsMounted && !state.ShouldAoE && state.ActiveSeal != "Seal of Truth" && !state.OnGlobalCooldown {
			c.State.LastZealCast = now
			c.CastSpell("Seal of Truth")

			return
		}

		if state.IsAlive &&
			!state.IsMounted && state.ShouldAoE && state.ActiveSeal != "Seal of Righteousness" && !state.OnGlobalCooldown {
			c.State.LastZealCast = now
			c.CastSpell("Seal of Righteousness")
			return
		}
	}

	if state.IsAlive && !state.IsMounted && !state.BlessingOn && !state.InCombat {
		c.CastSpell("Blessing of Might")
		return
	}

	if state.IsAlive && !state.IsMounted && state.InCombat && state.IsJudgementReady && state.ActiveSeal != "none" && !state.JudgementsOfThePureActive {
		c.CastSpell("Judgement")
		return
	}

	if state.IsAlive && !state.IsMounted && state.InCombat && (c.PopCooldowns || c.ForceCooldowns) && state.ZealotryAvailable {
		c.CastSpell("Zealotry")
		return
	}

	if state.IsAlive && !state.OnGlobalCooldown && !state.IsMounted && state.InCombat && !state.InquisitionActive && (state.HolyPower > 0 || state.DivinePurposeActive) {
		c.CastSpell("Inquisition")
		return
	}

	if state.IsAlive && !state.IsMounted && state.InCombat && (c.PopCooldowns || c.ForceCooldowns) && state.ZealotryActive && state.AvengingWrathAvailable {
		c.CastSpell("Avenging Wrath", 10)
		c.CastSpell("Trinket 1")
		c.PopCooldowns = false
		return
	}

	if state.IsAlive && !state.IsMounted && state.InCombat && state.CrusaderStrikeAvailable && state.HolyPower != 3 {
		if state.ShouldAoE && state.DivineStormAvailable {
			c.CastSpell("Divine Storm")
			return
		}

		c.CastSpell("Crusader Strike")
		return
	}

	if state.IsAlive && !state.IsMounted && state.InCombat && state.TemplarsVerdictReady && (state.HolyPower == 3 || state.DivinePurposeActive) {
		c.CastSpell("Templar's Verdict")
		return
	}

	if state.IsAlive && !state.IsMounted && state.InCombat && state.ArtOfWarActive {
		c.CastSpell("Exorcism")
		return
	}

	if state.IsAlive && !state.IsMounted && state.InCombat && state.HammerOfWrathAvailable {
		c.CastSpell("Hammer of Wrath")
		return
	}

	if state.IsAlive && !state.IsMounted && state.InCombat && state.ConsecrationAvailable && state.ShouldAoE {
		c.CastSpell("Consecration")
		return
	}

	if state.IsAlive && !state.IsMounted && state.InCombat && state.IsJudgementReady && state.ActiveSeal != "none" {
		c.CastSpell("Judgement")
		return
	}

}

func (c *Retribution) UpdateTables() {
	// optionValues := c.TViewTableValues["options"]
	stateValues := c.TViewTableValues["state"]

	stateValues["Inquisition Active"] = classes.TableCellValue{ZIndex: 31, NameColor: tcell.ColorWhite, Value: fmt.Sprintf("%t", c.State.InquisitionActive), ValueColor: util.GetColor(c.State.InquisitionActive, tcell.ColorGreen, tcell.ColorRed)}

	c.Paladin.UpdateTables()
}
