# In Mouth Probe (IMP)

## Features

- **Ingenious:** Utilizes temperature and humidity changes to predict device placement.
- **Stable:** Backend ensures efficient and precise data processing.
- **Fresh:** Minimalistic frontend UI for a seamless user experience.

## Tech Specs (Speculations)

### Hardware:
  - Based on the [WM1110 module dev-kit](https://wiki.seeedstudio.com/wio_tracker_for_sidewalk/)
    - **Chips:** LR1110, nRF52840
    - **Sensors:**
      - LIS3DH Triple-Axis Accelerometer
      - Sensirion SHT41-AD1F-R2 Humidity Sensor

### Software:

<img src="./.github/img/system-diagram.jpg" alt="system diagram" width="400"/>

  - **Firmware:** Zephyr RTOS on the WM1110A dev kit for sidewalk
  - **Middleware:** Golang, using [this](go.bug.st/serial) serial library
  - **Frontend:** Svelte5

## Data Model and Functionality

### Sensors and Prediction:
- **Input Data:**
  - Onboard temperature and humidity sensor
  - 3-axis accelerometer
- **Processing:**
  - Baseline data created by the user when the device is not in the mouth.
  - Data is refreshed every 300ms.
  - An algorithm analyzes trends in humidity and temperature changes, applying a derivative function for faster detection.
  - The movement of the data is also considered

### Display:
<img src="./.github/img/imp-ui.jpg" alt="ui screenshot" width="400"/>
