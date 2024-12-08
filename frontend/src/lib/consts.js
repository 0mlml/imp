export const DATA_STATE_DEVICE_UNKNOWN = 0;
export const DATA_STATE_DEVICE_NOT_RESPONDING = 1;
export const DATA_STATE_LOADING = 2;
export const DATA_STATE_IN_MOUTH = 3;
export const DATA_STATE_OUT_MOUTH = 4;

export const TREND_WINDOW_SIZE = 15;
export const TEMP_THRESHOLD = 1.2;
export const HUMIDITY_THRESHOLD = 70;
export const MIN_INCREASING_SAMPLES = 5;
export const MIN_DECREASING_SAMPLES = 3;
export const DERIVATIVE_THRESHOLD = 0.018; 

export const MINIMUM_TEMPERATURE_THRESHOLD = 0.35;
export const MINIMUM_HUMIDITY_THRESHOLD = 12;

export const ACCELERATION_DERIVATIVE_THRESHOLD = 0.2;
export const MIN_STABLE_SAMPLES = 8;

export const API_BASE = 'http://localhost:8080/api';