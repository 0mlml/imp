package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
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
				if strings.Contains(line, "SERIALPUSH>") {
					sensorData, err := parseSerialData(line)
					if err != nil {
						fmt.Printf("Parse error: %v\n", err)
						continue
					}

					// fmt.Printf("Temperature: %.2f°C, Humidity: %.2f%%\n",
					// 	sensorData.Humidity.Temperature,
					// 	sensorData.Humidity.Humidity)
					// fmt.Printf("Acceleration - X: %.2f, Y: %.2f, Z: %.2f, Peak: %.2f\n",
					// 	sensorData.Accelerometer.X,
					// 	sensorData.Accelerometer.Y,
					// 	sensorData.Accelerometer.Z,
					// 	sensorData.Accelerometer.PeakAcceleration)
					addToSensorBacklog(*sensorData)
				} else {
					// fmt.Print(line)
				}
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

var logRegex = regexp.MustCompile(`\[\d{2}:\d{2}:\d{2}\.\d{3},\d{3}\] <[a-z]+> [^:]+: SERIALPUSH>(.+)<SERIALPUSH`)

func parseSerialData(logLine string) (*SensorData, error) {
	// Try simpler regex first to just get the data between SERIALPUSH markers
	pushRegex := regexp.MustCompile(`SERIALPUSH>([^<]+)`)
	matches := pushRegex.FindStringSubmatch(logLine)
	if len(matches) != 2 {
		return nil, fmt.Errorf("no SERIALPUSH data found in: %s", logLine)
	}

	dataStr := matches[1]
	pairs := strings.Split(dataStr, ",")

	data := &SensorData{
		Humidity:      &HumidityData{},
		Accelerometer: &AccelerometerData{},
	}

	for _, pair := range pairs {
		kv := strings.Split(pair, ":")
		if len(kv) != 2 {
			return nil, fmt.Errorf("invalid key-value pair: %s", pair)
		}

		key := strings.TrimSpace(kv[0])
		value, err := strconv.ParseFloat(strings.TrimSpace(kv[1]), 64)
		if err != nil {
			return nil, fmt.Errorf("error parsing value for %s: %v", key, err)
		}

		switch key {
		case "temp":
			data.Humidity.Temperature = value
		case "hum":
			data.Humidity.Humidity = value
		case "motx":
			data.Accelerometer.X = value
		case "moty":
			data.Accelerometer.Y = value
		case "motz":
			data.Accelerometer.Z = value
		case "peak":
			data.Accelerometer.PeakAcceleration = value
		default:
			return nil, fmt.Errorf("unknown key: %s", key)
		}
	}

	return data, nil
}
