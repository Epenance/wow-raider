package paladin

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
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
	c.State.ActiveSeal = "none"

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

func (c *Paladin) UpdateTables() {
	stateValues := c.TViewTableValues["state"]

	/*
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
	*/

	stateValues["Crusader Strike Available"] = classes.TableCellValue{ZIndex: 31, NameColor: tcell.ColorWhite, Value: fmt.Sprintf("%t", c.State.CrusaderStrikeAvailable), ValueColor: util.GetColor(c.State.CrusaderStrikeAvailable, tcell.ColorGreen, tcell.ColorRed)}
	stateValues["Judgement Ready"] = classes.TableCellValue{ZIndex: 32, NameColor: tcell.ColorWhite, Value: fmt.Sprintf("%t", c.State.IsJudgementReady), ValueColor: util.GetColor(c.State.IsJudgementReady, tcell.ColorGreen, tcell.ColorRed)}
	stateValues["Active Seal"] = classes.TableCellValue{ZIndex: 33, NameColor: tcell.ColorWhite, Value: c.State.ActiveSeal, ValueColor: tcell.ColorWhite}

	stateValues["Should AoE"] = classes.TableCellValue{ZIndex: 40, NameColor: tcell.ColorWhite, Value: fmt.Sprintf("%t", c.State.ShouldAoE), ValueColor: util.GetColor(c.State.ShouldAoE, tcell.ColorGreen, tcell.ColorRed)}
	stateValues["Consecration Available"] = classes.TableCellValue{ZIndex: 41, NameColor: tcell.ColorWhite, Value: fmt.Sprintf("%t", c.State.ConsecrationAvailable), ValueColor: util.GetColor(c.State.ConsecrationAvailable, tcell.ColorGreen, tcell.ColorRed)}
	stateValues["Blessing On"] = classes.TableCellValue{ZIndex: 42, NameColor: tcell.ColorWhite, Value: fmt.Sprintf("%t", c.State.BlessingOn), ValueColor: util.GetColor(c.State.BlessingOn, tcell.ColorGreen, tcell.ColorRed)}
	stateValues["Hammer of Wrath Available"] = classes.TableCellValue{ZIndex: 43, NameColor: tcell.ColorWhite, Value: fmt.Sprintf("%t", c.State.HammerOfWrathAvailable), ValueColor: util.GetColor(c.State.HammerOfWrathAvailable, tcell.ColorGreen, tcell.ColorRed)}
	stateValues["Avenging Wrath Available"] = classes.TableCellValue{ZIndex: 44, NameColor: tcell.ColorWhite, Value: fmt.Sprintf("%t", c.State.AvengingWrathAvailable), ValueColor: util.GetColor(c.State.AvengingWrathAvailable, tcell.ColorGreen, tcell.ColorRed)}
	stateValues["Avenging Wrath Active"] = classes.TableCellValue{ZIndex: 45, NameColor: tcell.ColorWhite, Value: fmt.Sprintf("%t", c.State.AvengingWrathActive), ValueColor: util.GetColor(c.State.AvengingWrathActive, tcell.ColorGreen, tcell.ColorRed)}

	c.BaseClass.UpdateTables()
}
