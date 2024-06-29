package shaman

import (
	"wow-raider/classes"
	"wow-raider/util"
)

type ShamanState struct {
	classes.BaseState
	PrimalStrikeAvailable  bool
	FlameShockAvailable    bool
	FlameShockDotActive    bool
	EarthShockAvailable    bool
	WindfuryMissing        bool
	FlametongueMissing     bool
	LightningShieldMissing bool
}

type Shaman struct {
	classes.BaseClass
	State ShamanState
}

func (c *Shaman) Init(listeners []classes.KeyListener) error {
	c.Class = "Shaman"

	if err := c.BaseClass.Init(listeners); err != nil {
		return err
	}

	return nil
}

func (c *Shaman) SetState() {
	c.BaseClass.SetState()

	c.SyncState(&c.BaseClass.State, &c.State)

	c.State.PrimalStrikeAvailable = c.CheckColor(util.BLUE, 0, 0)
	c.State.FlameShockAvailable = c.CheckColor(util.BLUE, 5, 0)
	c.State.FlameShockDotActive = c.CheckColor(util.GREEN, 30, 0)
	c.State.EarthShockAvailable = c.CheckColor(util.BLUE, 15, 0)
	c.State.WindfuryMissing = c.CheckColor(util.RED, 20, 5)
	c.State.FlametongueMissing = c.CheckColor(util.RED, 25, 5)
	c.State.LightningShieldMissing = c.CheckColor(util.RED, 35, 5)
}

func (c *Shaman) UpdateTables() {
	// stateValues := c.TViewTableValues["state"]

	c.BaseClass.UpdateTables()
}
