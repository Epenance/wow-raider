package shaman

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"wow-raider/classes"
	"wow-raider/util"
)

type ShamanState struct {
	classes.BaseState
	ChainLightningAvailable bool
	PrimalStrikeAvailable   bool
	FlameShockAvailable     bool
	FlameShockDotActive     bool
	EarthShockAvailable     bool
	WindfuryMissing         bool
	FlametongueMissing      bool
	LightningShieldMissing  bool
	FlameshockDotsActive    int
	FireNovaAvailable       bool
	ShouldAoE               bool
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
	c.State.PrimalStrikeAvailable = c.CheckColor(util.BLUE, 45, 5)
	c.State.FlameShockAvailable = c.CheckColor(util.BLUE, 5, 0)
	c.State.FlameShockDotActive = c.CheckColor(util.GREEN, 30, 0)
	c.State.ShouldAoE = c.CheckColor(util.GREEN, 50, 5)
	c.State.EarthShockAvailable = c.CheckColor(util.BLUE, 15, 0)
	c.State.WindfuryMissing = c.CheckColor(util.RED, 20, 5)
	c.State.FlametongueMissing = c.CheckColor(util.RED, 25, 5)
	c.State.LightningShieldMissing = c.CheckColor(util.RED, 35, 5)
	c.State.FireNovaAvailable = c.CheckColor(util.BLUE, 45, 0)

	flameShockDots3 := c.CheckColor(util.GREEN, 40, 0)

	if flameShockDots3 {
		c.State.FlameshockDotsActive = 3
	} else {
		c.State.FlameshockDotsActive = 0
	}

}

func (c *Shaman) UpdateTables() {
	stateValues := c.TViewTableValues["state"]
	stateValues["Should AoE"] = classes.TableCellValue{ZIndex: 1, NameColor: tcell.ColorWhite, Value: fmt.Sprintf("%t", c.State.ShouldAoE), ValueColor: util.GetColor(c.State.ShouldAoE, tcell.ColorGreen, tcell.ColorRed)}

	c.BaseClass.UpdateTables()
}
