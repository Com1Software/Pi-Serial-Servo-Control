package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/googolgl/go-i2c"
	"github.com/googolgl/go-pca9685"
	"go.bug.st/serial"
)

func main() {
	fmt.Println("Land Vehicle Test System")

	i2c, err := i2c.New(pca9685.Address, "/dev/i2c-1")
	if err != nil {
		log.Fatal(err)
	}
	pca0, err0 := pca9685.New(i2c, nil)
	if err0 != nil {
		log.Fatal(err0)
	}
	pca1, err1 := pca9685.New(i2c, nil)
	if err1 != nil {
		log.Fatal(err1)
	}
	pca1.SetChannel(1, 0, 130)
	pca0.SetChannel(0, 0, 130)
	servo1 := pca1.ServoNew(1, nil)
	servo0 := pca0.ServoNew(0, nil)

	// Retrieve the port list
	ports, err := serial.GetPortsList()
	if err != nil {
		log.Fatal(err)
	}
	if len(ports) == 0 {
		log.Fatal("No serial ports found!")
	}

	// Print the list of detected ports
	for _, port := range ports {
		fmt.Printf("Found port: %v\n", port)
	}

	// Open the first serial port detected at 9600bps N81
	mode := &serial.Mode{
		BaudRate: 115200,
		Parity:   serial.NoParity,
		DataBits: 8,
		StopBits: serial.OneStopBit,
	}
	port, err := serial.Open(ports[0], mode)
	if err != nil {
		log.Fatal(err)
	}

	for {

		line := ""
		buff := make([]byte, 1)
		on := true
		for on != false {
			line = ""
			for {
				n, err := port.Read(buff)
				if err != nil {
					log.Fatal(err)
				}
				if n == 0 {
					fmt.Println("\nEOF")
					break
				}
				line = line + string(buff[:n])
				if strings.Contains(string(buff[:n]), "\n") {
					break
				}

			}

			ch1, ch2, ch3, ch4 := getCHPosition(line)
			fmt.Print("\033[u\033[K")
			fmt.Println("Land Vehicle Test System")
			fmt.Printf("CH1=%s CH2=%s CH3=%s CH4=%s\n", ch1, ch2, ch3, ch4)
			fmt.Println("======================================")
			sv, _ := strconv.Atoi(ch1)
			svx := sv / 23
			sv1 := 1 + svx
			sv2 := sv1
			if sv1 > 65 {
				if sv1 > 70 {
					sv2 = sv2 + 5
					if sv1 > 75 {
						sv2 = sv2 + 5
						if sv1 > 85 {
							sv2 = sv2 + 5
						}
					}
				}
			} else {
				if sv1 < 60 {
					sv2 = sv2 - 5
					if sv1 < 55 {
						sv2 = sv2 - 5
						if sv1 < 45 {
							sv2 = sv2 - 5
						}
					}
				}

			}

			fmt.Printf("Servo Position: %d   %d   %d \n", svx, sv1, sv2)
			servo1.Angle(sv2)

			servo0.Angle(50)

		}
	}
}

func getCHPosition(sentence string) (string, string, string, string) {
	data := strings.Split(sentence, ",")
	ch1 := ""
	ch2 := ""
	ch3 := ""
	ch4 := ""
	if len(data) == 4 {
		ch1data := strings.Split(data[0], "=")
		ch2data := strings.Split(data[1], "=")
		ch3data := strings.Split(data[2], "=")
		ch4data := strings.Split(data[3], "=")

		if string(ch1data[0]) == "CH1" {
			ch1 = ch1data[1]
		}
		if string(ch2data[0]) == "CH2" {
			ch2 = ch2data[1]
		}
		if string(ch3data[0]) == "CH3" {
			ch3 = ch3data[1]
		}
		if string(ch4data[0]) == "CH4" {
			ch4 = ch4data[1]
		}
	}
	return ch1, ch2, ch3, ch4

}
