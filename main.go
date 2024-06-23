package main

import (
	"wow-raider/classes/paladin"
)

func main() {
	routine := &paladin.Retribution{}

	err := routine.Init()
	defer routine.Uninit()

	if err != nil {
		return
	}

	routine.Run()
}
