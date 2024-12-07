export const DATA_STATE_DEVICE_UNKNOWN = 0;
export const DATA_STATE_DEVICE_NOT_RESPONDING = 1;
export const DATA_STATE_LOADING = 2;
export const DATA_STATE_IN_MOUTH = 3;
export const DATA_STATE_OUT_MOUTH = 4;

export const THRESHOLD_TEMP_IN_MOUTH = 27;  // Temperature in °C for the mouth
export const THRESHOLD_TEMP_DIFFERENCE = 2; // Maximum allowed difference between environmental and mouth temperature
export const THRESHOLD_HUMIDITY_IN_MOUTH = 80; // Humidity in % for the mouth
export const THRESHOLD_HUMIDITY_DIFFERENCE = 20; // Maximum allowed difference between environmental and mouth humidity

export const THRESHOLD_STRICT_ACCEL = 0.1; // Maximum allowed movement for sync (in m/s²)
export const THRESHOLD_ACCEL = 0.2; // Maximum allowed movement for mouth (in m/s²)
export const THRESHOLD_STILL_PERIOD = 3000/300; // Minimum still period (in updates; milliseconds/update_rate)
export const GRAVITY_ACCEL = 9.9; // Gravity acceleration (in m/s²)
export const TEMP_MAX_DELTA = 3; // Maximum acceptable temperature change between readings
export const HUM_MAX_DELTA = 7; // Maximum acceptable humidity change between readings

export const API_BASE = 'http://localhost:8080/api';