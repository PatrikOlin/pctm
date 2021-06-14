package main

import (
	"bufio"
	"log"
	"machine"
	"os"
	"time"

	"github.com/alphahorizonio/tinynet/pkg/tinynet"
)

var (
	BUFLEN = 1038
)

const BASEURL = "http://bl-kiosk-map.jls-sto3.elastx.net/map/"

func main() {
	button := machine.D5
	button.Configure(machine.PinConfig{Mode: machine.PinInput})

	for {
		if button.Get() {
			on()
		}
		off()
	}
}

func on() {
	led := machine.LED
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})

	led.High
	time.Sleep(time.Millisecond * 50)
	checkInOrOut("signIn")
}

func off() {
	led := machine.LED
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})

	led.Low
	time.Sleep(time.Millisecond * 50)
	checkInOrOut("signOut")
}

func checkInOrOut(endpoint string) {
	userId := "2"
	conn, err := tinynet.Dial("tcp", BASEURL+endpoint+"/"+userId)
	// resp, err := http.Post(BASEURL+endpoint+"/"+userId, "application/json", nil)
	if err != nil {
		log.Fatalf("oops, error! %v", err)
	}

	reader := bufio.NewReader(os.Stdin)

	for {
		out, err := reader.ReadString('\n')
		if err != nil {
			log.Println("could not read from stdin", err)

			os.Exit(1)
		}

		if n, err := conn.Write([]byte(out)); err != nil {
			if n == 0 {
				break
			}

			log.Println("could not write from connection, removing connection", err)

			break
		}

		buf := make([]byte, BUFLEN)
		if n, err := conn.Read(buf); err != nil {
			if n == 0 {
				break
			}

			log.Println("could not read from connection, removing connectoin", err)

			break
		}

		log.Println(string(buf))
	}

	log.Println("Disconnected")

	if err := conn.Close(); err != nil {
		log.Println("could not close connection", err)
	}
}
