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
			err := c.CaptureGame()

			if err != nil {
				return
			}

			c.SetState()
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
}

func (c *Retribution) Rotation() {
	// Do stuff
}
