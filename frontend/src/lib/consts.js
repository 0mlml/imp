export const DATA_STATE_DEVICE_UNKNOWN = 0;
export const DATA_STATE_DEVICE_NOT_RESPONDING = 1;
export const DATA_STATE_LOADING = 2;
export const DATA_STATE_IN_MOUTH = 3;
export const DATA_STATE_OUT_MOUTH = 4;

export const TREND_WINDOW_SIZE = 12; // Shorter window
export const TEMP_THRESHOLD = 1.2; // Lower temperature threshold
export const HUMIDITY_THRESHOLD = 70; // Lower humidity threshold
export const MIN_INCREASING_SAMPLES = 5; // Fewer samples needed for trend
export const MIN_DECREASING_SAMPLES = 5;
export const DERIVATIVE_THRESHOLD = 0.018; 

export const MINIMUM_TEMPERATURE_THRESHOLD = 0.35;
export const MINIMUM_HUMIDITY_THRESHOLD = 12;

export const API_BASE = 'http://localhost:8080/api';