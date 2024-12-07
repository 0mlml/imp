package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"sync"
	"syscall"
	"time"

	"go.bug.st/serial"
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

type HumidityData struct {
	Temperature float64 `json:"temperature"`
	Humidity    float64 `json:"humidity"`
}

type AccelerometerData struct {
	X                float64 `json:"x"`
	Y                float64 `json:"y"`
	Z                float64 `json:"z"`
	PeakAcceleration float64 `json:"peak_acceleration"`
}

type SensorData struct {
	Humidity      *HumidityData      `json:"humidity"`
	Accelerometer *AccelerometerData `json:"accelerometer"`
	BatteryLevel  float64            `json:"battery_level"`
}
