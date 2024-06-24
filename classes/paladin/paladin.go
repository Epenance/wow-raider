package paladin

import (
	"reflect"
	"time"
	"wow-raider/classes"
	"wow-raider/util"
)

type PaladinState struct {
	classes.BaseState
	HolyPower               int
	CrusaderStrikeAvailable bool
	IsJudgementReady        bool
	ActiveSeal              string
	ShouldAoE               bool
	ConsecrationAvailable   bool
	BlessingOn              bool
	HammerOfWrathAvailable  bool
	AvengingWrathAvailable  bool
	AvengingWrathActive     bool
	LastZealCast            time.Time
}

type Paladin struct {
	classes.BaseClass
	State PaladinState
}

func (c *Paladin) PrintState() {
	// Use reflect.ValueOf(c).Elem() to get the correct value
	classes.PrintFields(reflect.ValueOf(c.State))
}

func (c *Paladin) Init() error {
	c.Class = "Paladin"

	if err := c.BaseClass.Init(); err != nil {
		return err
	}

	return nil
}

func (c *Paladin) SetState() {
	c.BaseClass.SetState()

	// Hack because there is no inheritance in Go
	c.SyncState(&c.BaseClass.State, &c.State)

	// Paladin specific state
	c.State.CrusaderStrikeAvailable = c.CheckColor(util.BLUE, 25, 0)
	c.State.IsJudgementReady = c.CheckColor(util.BLUE, 5, 0)
	c.State.ShouldAoE = c.CheckColor(util.GREEN, 35, 0)
	c.State.ConsecrationAvailable = c.CheckColor(util.YELLOW, 45, 0)
	c.State.BlessingOn = !c.CheckColor(util.RED, 25, 5)
	c.State.HammerOfWrathAvailable = c.CheckColor(util.BLUE, 20, 0)
	c.State.AvengingWrathAvailable = c.CheckColor(util.BLUE, 65, 5)
	c.State.AvengingWrathActive = c.CheckColor(util.GREEN, 60, 5)

	c.SetZeal()
	c.SetHolyPower()
}

func (c *Paladin) SetZeal() {
	noActiveSeal := c.CheckColor(util.RED, 0, 0)

	if noActiveSeal {
		c.State.ActiveSeal = "none"
		return
	}

	activeSealIsRighteousness := c.CheckColor(util.GREEN, 50, 0)

	if activeSealIsRighteousness {
		c.State.ActiveSeal = "Seal of Righteousness"
		return
	}

	activeSealIsTruth := c.CheckColor(util.BLUE, 50, 0)

	if activeSealIsTruth {
		c.State.ActiveSeal = "Seal of Truth"
		return
	}
}

func (c *Paladin) SetHolyPower() {
	isHolyPower1 := c.CheckColor(util.GREEN, 45, 5)
	isHolyPower2 := c.CheckColor(util.BLUE, 45, 5)
	isHolyPower3 := c.CheckColor(util.YELLOW, 45, 5)

	if isHolyPower1 {
		c.State.HolyPower = 1
		return
	}

	if isHolyPower2 {
		c.State.HolyPower = 2
		return
	}

	if isHolyPower3 {
		c.State.HolyPower = 3
		return
	}

	c.State.HolyPower = 0
}
