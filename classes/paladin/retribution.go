package paladin

import (
	"wow-raider/util"
)

type RetributionState struct {
	PaladinState
}

type Retribution struct {
	Paladin
	State RetributionState
}

func (c *Retribution) Run() {
	for !c.InterruptProgram {
		if c.RunProgram {
			// Do stuff
		}
	}
}

func (c *Retribution) Setup() {
	util.Log("Setting up Retribution routine")
	c.Paladin.Setup()
}
