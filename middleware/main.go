package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"strings"

	"go.bug.st/serial"
	"go.bug.st/serial/enumerator"
)

func startConsoleInput(serialMgr *SerialManager) {
	scanner := bufio.NewScanner(os.Stdin)
	go func() {
		for scanner.Scan() {
			input := scanner.Text()
			if err := serialMgr.WriteData([]byte(input)); err != nil {
				fmt.Printf("Failed to write: %v\n", err)
			}
		}
	}()
}

const (
	SENSOR_BACKLOG_SIZE = 100
)

var (
	debugFlag     = flag.Bool("debug", false, "Enable debug mode")
	serialNumber  = flag.String("serial", "", "Serial number of the WM1110 board")
	portName      = flag.String("port", "", "Name of the serial port to connect to")
	sensorBacklog = make([]string, 0)
)

func addToSensorBacklog(sd SensorData) {
	if len(sensorBacklog) >= SENSOR_BACKLOG_SIZE {
		sensorBacklog = sensorBacklog[1:]
	}
	jsonData, err := json.Marshal(sd)
	if err != nil {
		fmt.Printf("Failed to marshal sensor data: %v\n", err)
		return
	}
	sensorBacklog = append(sensorBacklog, string(jsonData))
}

func addFakeData() {
	sd := SensorData{
		Humidity: &HumidityData{
			Humidity:    rand.Float64()*50 + 50,
			Temperature: rand.Float64()*10 + 20,
		},
		Accelerometer: &AccelerometerData{
			X:                rand.Float64()*2 - 1,
			Y:                rand.Float64()*2 - 1,
			Z:                rand.Float64()*2 - 1,
			PeakAcceleration: rand.Float64()*10 + 10,
		},
		BatteryLevel: rand.Float64()*20 + 80,
	}
	addToSensorBacklog(sd)
}

func enumeratePorts() {
	ports, err := enumerator.GetDetailedPortsList()
	if err != nil {
		fmt.Printf("error getting port list: %v", err)
		return
	}

	fmt.Printf("Found %d ports:\n", len(ports))
	for _, port := range ports {
		fmt.Printf("  Name: %s\n", port.Name)
		fmt.Printf("  SerialNumber: %s\n", port.SerialNumber)
		fmt.Printf("  VID: %s\n", port.VID)
		fmt.Printf("  PID: %s\n", port.PID)
		fmt.Printf("  Product: %s\n", port.Product)
		fmt.Printf("  IsUSB: %t\n", port.IsUSB)
	}
}

func findWM1110Port() (string, error) {
	ports, err := enumerator.GetDetailedPortsList()
	if err != nil {
		return "", fmt.Errorf("error getting port list: %v", err)
	}

	for _, port := range ports {
		if port.IsUSB && strings.Contains(port.SerialNumber, *serialNumber) {
			return port.Name, nil
		}
	}

	return "", fmt.Errorf("WM1110 board not found")
}

func openSerialConnection(portName string) (serial.Port, error) {
	mode := &serial.Mode{
		BaudRate: 115200,
	}

	port, err := serial.Open(portName, mode)
	if err != nil {
		return nil, fmt.Errorf("error opening port %s: %v", portName, err)
	}

	return port, nil
}

func debugMode() {
	for i := 0; i < SENSOR_BACKLOG_SIZE; i++ {
		addFakeData()
	}
}

func getSerialPort() (serial.Port, error) {
	if *portName == "" && *serialNumber == "" {
		fmt.Println("No serial number or port name specified, enumerating ports...")
		enumeratePorts()
		return nil, fmt.Errorf("no serial number or port name specified")
	}

	var port string
	var err error
	if *serialNumber != "" {
		fmt.Println("Enumerating ports to find WM1110 board...")
		port, err = findWM1110Port()
		fmt.Printf("Found WM1110 board at port %s\n", port)
		if err != nil {
			fmt.Printf("Failed to find WM1110 board: %v\n", err)
			return nil, err
		}
	}
	if *portName != "" && port == "" {
		fmt.Printf("Using port %s\n", *portName)
		port = *portName
	}
	if port == "" {
		fmt.Println("No port found")
		return nil, errors.New("no port found")
	}

	sPort, err := openSerialConnection(port)
	if err != nil {
		fmt.Printf("Failed to open serial connection: %v\n", err)
		return nil, err
	}

	sPort.SetMode(&serial.Mode{
		BaudRate: 115200,
		DataBits: 8,
		Parity:   serial.NoParity,
		StopBits: serial.OneStopBit,
	})

	return sPort, nil
}

func main() {
	flag.Parse()

	if *debugFlag {
		debugMode()
		fmt.Printf("Debug mode enabled. Generated %d fake sensor data entries.\n", SENSOR_BACKLOG_SIZE)
		fmt.Println("There will be no serial communication.")
	} else {
		go func() {
			sPort, err := getSerialPort()
			if err != nil {
				os.Exit(1)
			}
			defer sPort.Close()
			fmt.Println("Successfully connected to WM1110!")

			serialMgr := NewSerialManager(sPort)
			serialMgr.StartReadLoop()

			startConsoleInput(serialMgr)

			fmt.Println("Serial communication started. Type your messages and press Enter to send.")
			fmt.Println("Press Ctrl+C to exit.")

			sigChan := make(chan os.Signal, 1)
			signal.Notify(sigChan, os.Interrupt)

			<-sigChan
			fmt.Println("\nShutting down...")
			serialMgr.Stop()
		}()
	}

	startHttpServer()
}
