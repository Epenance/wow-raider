package paladin

import (
	"fmt"
	"wow-raider/util"
)

type RetributionState struct {
	PaladinState
	DivineStormAvailable bool
	ArtOfWarActive       bool
	TemplarsVerdictReady bool
	DivinePurposeActive  bool
	InquisitionActive    bool
	ZealotryActive       bool
	ZealotryAvailable    bool
}

type Retribution struct {
	Paladin
	State RetributionState
}

func (c *Retribution) Run() {
	for !c.InterruptProgram {
		if c.RunProgram {
			// Do stuff
			err := c.CaptureGame()

			if err != nil {
				return
			}

			c.SetState()

			fmt.Println(c.State.IsMounted)
		}
	}
}

func (c *Retribution) Init() error {
	c.Spec = "Retribution"

	if err := c.Paladin.Init(); err != nil {
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
}

func (c *Retribution) Rotation() {
	// Do stuff
}
