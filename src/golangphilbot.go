package main

import (
	"fmt"
	"time"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/sphero"
)

func main() {
	gbot := gobot.NewGobot()

	adaptor := sphero.NewSpheroAdaptor("Sphero", "/dev/tty.Sphero-OOR-AMP-SPP")
	spheroDriver := sphero.NewSpheroDriver(adaptor, "sphero")
	spheroDriver.SetRGB(255, 200, 100)
	color := true
	work := func() {
		gobot.Every(2*time.Second, func() {
			color = !color
			if color {
				spheroDriver.SetRGB(200, 200, 0)
				fmt.Printf("200,200,0\n")
			} else {
				spheroDriver.SetRGB(20, 200, 100)
				fmt.Printf("20,200,100\n")
			}
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
