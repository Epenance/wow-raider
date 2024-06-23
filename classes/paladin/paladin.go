package paladin

import (
	"reflect"
	"time"
	"wow-raider/classes"
)

type PaladinState struct {
	classes.BaseState
	HolyPower               int
	CrusaderStrikeAvailable bool
	IsJudgementReady        bool
	ActiveSeal              string
	ShouldAoE               bool
	ConsecrationAvailable   bool
	BlesssingOn             bool
	HammerOfWrathAvailable  bool
	AvengeWrathAvailable    bool
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
