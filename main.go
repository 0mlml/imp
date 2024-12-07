package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"go.bug.st/serial"
	"go.bug.st/serial/enumerator"
)

type SerialManager struct {
	port     serial.Port
	stopChan chan struct{}
	wg       sync.WaitGroup
}

func NewSerialManager(port serial.Port) *SerialManager {
	return &SerialManager{
		port:     port,
		stopChan: make(chan struct{}),
	}
}

func (sm *SerialManager) StartReadLoop() {
	sm.wg.Add(1)
	go func() {
		defer sm.wg.Done()

		reader := bufio.NewReader(sm.port)

		for {
			select {
			case <-sm.stopChan:
				return
			default:
				line, err := reader.ReadString('\n')
				if err != nil {
					if err != io.EOF && err != syscall.ETIMEDOUT {
						fmt.Printf("Read error: %v\n", err)
					}
					time.Sleep(10 * time.Millisecond)
					continue
				}

				line = string(bytes.ReplaceAll([]byte(line), []byte("\r"), []byte("")))
				fmt.Print(line)
			}
		}
	}()
}

func (sm *SerialManager) WriteData(data []byte) error {
	if !bytes.HasSuffix(data, []byte("\n")) {
		data = append(data, '\n')
	}

	_, err := sm.port.Write(data)
	if err != nil {
		return fmt.Errorf("write error: %v", err)
	}
	return nil
}

func (sm *SerialManager) Stop() {
	close(sm.stopChan)
	sm.wg.Wait()
}

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

var (
	serialNumber = flag.String("serial", "", "Serial number of the WM1110 board")
	portName     = flag.String("port", "", "Name of the serial port to connect to")
)

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

func main() {
	flag.Parse()
	if *portName == "" && *serialNumber == "" {
		fmt.Println("No serial number or port name specified, enumerating ports...")
		enumeratePorts()
		return
	}

	var port string
	var err error
	if *serialNumber != "" {
		fmt.Println("Enumerating ports to find WM1110 board...")
		port, err = findWM1110Port()
		fmt.Printf("Found WM1110 board at port %s\n", port)
		if err != nil {
			fmt.Printf("Failed to find WM1110 board: %v\n", err)
			return
		}
	}
	if *portName != "" && port == "" {
		fmt.Printf("Using port %s\n", *portName)
		port = *portName
	}
	if port == "" {
		fmt.Println("No port found")
		return
	}

	sPort, err := openSerialConnection(port)
	if err != nil {
		fmt.Printf("Failed to open serial connection: %v\n", err)
		return
	}
	defer sPort.Close()

	sPort.SetMode(&serial.Mode{
		BaudRate: 115200,
		DataBits: 8,
		Parity:   serial.NoParity,
		StopBits: serial.OneStopBit,
	})

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
}
