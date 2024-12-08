# In Mouth Probe (IMP)

## 📖 Introduction

Welcome to the **In Mouth Probe (IMP)** GitHub repository! This project was developed as part of the course **Fundamentals of Smart Systems (CT60A4800)** by **Group G**. Our aim was to explore innovative applications of smart systems, resulting in the creation of IMP—a device that humorously yet effectively solves the problem of identifying its placement inside the mouth.

## 🚀 Features

- **Ingenious:** Utilizes temperature and humidity changes to predict device placement.
- **Stable:** Backend ensures efficient and precise data processing.
- **Fresh:** Minimalistic frontend UI for a seamless user experience.

## 🛠️ Technical Specifications

![Overview Diagram](./system-diagram.jpg)

IMP is built on cutting-edge hardware and software:
- **Hardware:**
  - Based on the [WM1110 module dev-kit](https://wiki.seeedstudio.com/wio_tracker_for_sidewalk/)
  - **Chips:** LR1110, nRF52840
  - **Connectivity:** AWS Sidewalk, LoRa (2-10km), BLE (<10m), BMesh, Thread, Zigbee, 802.15.4, ANT, 2.4 GHz
  - **Sensors:**
    - LIS3DH Triple-Axis Accelerometer
    - Sensirion SHT41-AD1F-R2 Humidity Sensor
- **Software:**
  - **Middleware:** Zephyr RTOS
  - **Frontend:** Built with Svelte for a blazing-fast user interface

## ⚙️ How It Works

The **In Mouth Probe (IMP)** system leverages a combination of sensors, algorithms, and user input to detect its placement in or out of the mouth with precision:

1. **Initialization and Calibration:**
   - The user presses a button to establish baseline data when the device is not in the mouth.
   - This baseline captures ambient temperature, humidity, and movement data to set a reference point for detection.

2. **Continuous Monitoring:**
   - The device collects real-time data every 300ms from the onboard sensors:
     - **Humidity and Temperature:** Detects rapid changes indicative of the moist, warm environment inside the mouth.
     - **Accelerometer:** Confirms movement patterns consistent with mouth placement.

3. **Data Processing and Prediction:**
   - A specialized algorithm compares incoming sensor data to the baseline, identifying trends and changes.
   - A derivative function ensures snappy detection by amplifying noticeable shifts in environmental conditions.
   - Accelerometer data adds an additional layer of validation to ensure accurate placement detection.

4. **User Feedback:**
   - The device's user interface provides live updates on:
     - **Sensor Readings:** Displays real-time temperature, humidity, and movement data.
     - **Placement Status:** Indicates whether the probe is "In the mouth" or "Not in the mouth."

5. **Connectivity:**
   - BLE, LoRa, and AWS Sidewalk support ensure seamless data transfer and integration with external systems if needed.

This combination of sensors, intelligent algorithms, and responsive design makes the IMP a robust and innovative solution to its uniquely playful challenge.

## 📊 Data Model and Functionality

### Sensors and Prediction:
- **Input Data:**
  - Onboard temperature and humidity sensor
  - 3-axis accelerometer
- **Processing:**
  - Baseline data created by the user when the device is not in the mouth.
  - Data is refreshed every 300ms.
  - An algorithm analyzes trends in humidity and temperature changes, applying a derivative function for faster detection.
  - Movement data ensures accurate placement detection.

### Action and Display:
- The user interface provides live sensor readings and real-time detection status:
  - **States:** "In the mouth" or "Not in the mouth"

## 🌟 Why IMP?

### Smartness:
- Continuous monitoring using advanced accelerometer data and environmental sensors.
- Utilizes real-time algorithms to detect changes in humidity and temperature with high precision.

### Robust Design:
- Backed by modern middleware for reliability.
- Intuitive and stylish frontend design.
