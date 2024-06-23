package paladin

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
