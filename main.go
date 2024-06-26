package main

import (
	"fmt"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
	"github.com/manifoldco/promptui"
	"wow-raider/classes/paladin"
)

type Routine interface {
	Init() error
	Uninit()
	Run()
}

func main() {
	a := app.New()
	w := a.NewWindow("TODO App")

	w.SetContent(widget.NewLabel("TODOs will go here"))
	w.ShowAndRun()

	prompt := promptui.Select{
		Label: "Select Rotation",
		Items: []string{"Retribution Paladin", "Cancel"},
	}

	_, result, err := prompt.Run()
	if err != nil {
		fmt.Println("Error selecting option:", err)
		return
	}

	routines := make(map[string]Routine)
	routines["Retribution Paladin"] = &paladin.Retribution{}

	if result != "Cancel" {
		routine, exists := routines[result]
		if !exists {
			fmt.Println("No routine found for selected option")
			return
		}

		err = routine.Init()
		if err != nil {
			fmt.Println("Error initializing routine:", err)
			return
		}
		defer routine.Uninit()

		routine.Run()
	}
}
