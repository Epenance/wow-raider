package main

import (
	"wow-raider/classes/paladin"
	"wow-raider/util"
)

func main() {
	routine := paladin.Retribution{}

	err := routine.Init()
	if err != nil {
		util.Log("Failed to initialize routine")
		return
	}

	routine.Run()
}
