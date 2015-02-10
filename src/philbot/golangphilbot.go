package main

import (
	"fmt"
	"time"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/sphero"
)

func main() {
	gbot := gobot.NewGobot()

	adaptor := sphero.NewSpheroAdaptor("sphero", "/dev/tty.Sphero-OOR-AMP-SPP")
	spheroDriver := sphero.NewSpheroDriver(adaptor, "sphero")
	spheroDriver.SetRGB(255, 200, 100)
	color := true
	work := func() {
		// go in a square
		// spheroDriver.SetRGB(12, 245, 245)
		// spheroDriver.Roll(50, uint16(0))
		// time.Sleep(5 * time.Second)

		// spheroDriver.SetRGB(98, 12, 245)
		// spheroDriver.Roll(50, uint16(90))
		// time.Sleep(5 * time.Second)

		// spheroDriver.SetRGB(0, 0, 0)
		// spheroDriver.Roll(50, uint16(180))
		// time.Sleep(5 * time.Second)

		// spheroDriver.SetRGB(100, 28, 65)
		// spheroDriver.Roll(50, uint16(270))
		// time.Sleep(5 * time.Second)

		gobot.Every(1*time.Second, func() {
			color = !color
			if color {
				spheroDriver.SetRGB(200, 200, 0)
				fmt.Printf("200,200,0\n")
			} else {
				spheroDriver.SetRGB(20, 200, 100)
				fmt.Printf("20,200,100\n")
			}
			// direction := uint16(gobot.Rand(360))
			// fmt.Printf("direction=%v\n", direction)
			// spheroDriver.Roll(80, direction)
		})
	}

	robot := gobot.NewRobot("sphero",
		[]gobot.Connection{adaptor},
		[]gobot.Device{spheroDriver},
		work,
	)

	gbot.AddRobot(robot)

	gbot.Start()
}
