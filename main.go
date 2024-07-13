package main

import (
	"github.com/manifoldco/promptui"
	"wow-raider/classes/paladin"
	"wow-raider/classes/shaman"
)

type Routine interface {
	Init() error
	Uninit()
	Run()
}

func main() {
	prompt := promptui.Select{
		Label: "Select Rotation",
		Items: []string{"Retribution Paladin", "Enhancement Shaman", "Cancel"},
	}

	_, result, err := prompt.Run()

	routines := map[string]Routine{}

	routines["Retribution Paladin"] = &paladin.Retribution{}
	routines["Enhancement Shaman"] = &shaman.Enhancement{}

	if result != "Cancel" {
		routine := routines[result]
		err = routine.Init()
		defer routine.Uninit()

		if err != nil {
			return
		}

		routine.Run()
	}
}
