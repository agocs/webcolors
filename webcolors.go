package main

import (
	"log"
	"net/http"
	"time"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/firmata"
	"github.com/hybridgroup/gobot/platforms/gpio"
)

func main() {

	//initialize command channels

	colors := make(chan int, 500)

	go blink(colors)

	//initialize web server

	http.HandleFunc("/green/", func(w http.ResponseWriter, req *http.Request) {
		log.Println("green")
		colors <- 0
		http.Redirect(w, req, "/", 302)
		return
	})
	http.HandleFunc("/blue/", func(w http.ResponseWriter, req *http.Request) {
		log.Println("blue")
		colors <- 1
		http.Redirect(w, req, "/", 302)
		return
	})
	http.HandleFunc("/red/", func(w http.ResponseWriter, req *http.Request) {
		log.Println("red")
		colors <- 2
		http.Redirect(w, req, "/", 302)
		return
	})
	http.HandleFunc("/rainbow/", func(w http.ResponseWriter, req *http.Request) {
		log.Println("rainbow")
		colors <- 3
		http.Redirect(w, req, "/", 302)
		return
	})

	clientfs := http.FileServer(http.Dir("client"))
	http.Handle("/", clientfs)

	err := http.ListenAndServe(":9000", nil)
	if err != nil {
		close(colors)
		panic(err.Error())
	}
}

func blink(colors chan int) {
	//initialize robot friend

	gbot := gobot.NewGobot()

	firmataAdaptor := firmata.NewFirmataAdaptor("arduino", "/dev/cu.usbmodem1421")
	led_g := gpio.NewLedDriver(firmataAdaptor, "led", "13")
	led_r := gpio.NewLedDriver(firmataAdaptor, "led", "12")
	led_b := gpio.NewLedDriver(firmataAdaptor, "led", "11")
	led_rainbow := gpio.NewLedDriver(firmataAdaptor, "led", "10")

	work := func() {
		gobot.Every(1*time.Millisecond, func() {
			color, alive := <-colors

			if !alive {
				panic("Quittin'")
			}

			switch color {
			case 0:
				led_rainbow.Off()
				led_g.Toggle()
			case 1:
				led_rainbow.Off()
				led_b.Toggle()
			case 2:
				led_rainbow.Off()
				led_r.Toggle()
			case 3:
				led_g.Off()
				led_b.Off()
				led_r.Off()
				led_rainbow.On()
			}

		})
	}

	robot := gobot.NewRobot("bot",
		[]gobot.Connection{firmataAdaptor},
		[]gobot.Device{led_g},
		[]gobot.Device{led_r},
		[]gobot.Device{led_b},
		[]gobot.Device{led_rainbow},
		work,
	)

	gbot.AddRobot(robot)

	gbot.Start()
}
