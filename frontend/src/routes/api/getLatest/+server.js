import {
    DATA_STATE_DEVICE_UNKNOWN,
    DATA_STATE_DEVICE_NOT_RESPONDING,
    DATA_STATE_LOADING,
    DATA_STATE_IN_MOUTH,
    DATA_STATE_OUT_MOUTH,
    TREND_WINDOW_SIZE,
    TEMP_THRESHOLD,
    HUMIDITY_THRESHOLD,
    MIN_INCREASING_SAMPLES,
    MIN_DECREASING_SAMPLES,
    DERIVATIVE_THRESHOLD,
    MINIMUM_TEMPERATURE_THRESHOLD,
    MINIMUM_HUMIDITY_THRESHOLD,
    API_BASE
} from '$lib/consts';

export async function POST({ request }) {
    const { environmentalTemperature, environmentalHumidity } = await request.json();
    const data = await (await fetch(`http://localhost:8080/getlatest?count=${TREND_WINDOW_SIZE}`)).json();
    
    const { humidity, temperature } = data[0].humidity;
    const { x, y, z } = data[0].accelerometer;

    const tempHistory = data.map(d => d.humidity.temperature);
    const humidityHistory = data.map(d => d.humidity.humidity);

    const trends = analyzeTrends(tempHistory, humidityHistory);
    const thresholdsMet = checkThresholds(temperature, humidity, environmentalTemperature);

    // Get previous state
    let previousState = false;
    try {
        const prevData = await (await fetch(`http://localhost:8080/getlatest?count=1&offset=1`)).json();
        previousState = prevData[0]?.inMouth || false;
    } catch (e) {
        previousState = false;
    }

    const current = {temperature: temperature, humidity: humidity, environmentalTemperature: environmentalTemperature, environmentalHumidity: environmentalHumidity};

    // Determine state
    const inMouth = determineState(trends, thresholdsMet, previousState, current);

    return new Response(JSON.stringify({
        success: true,
        body: {
            environmentalHumidity,
            environmentalTemperature,
            temperature,
            humidity,
            x, y, z,
            inMouth,
            trends,
            thresholdsMet
        }
    }), {
        status: 200,
        headers: { 'Content-Type': 'application/json' },
    });
}

function analyzeTrends(tempHistory, humidityHistory) {
    let tempIncreasing = 0;
    let tempDecreasing = 0;
    let humidityIncreasing = 0;
    let humidityDecreasing = 0;
    
    for (let i = 1; i < tempHistory.length; i++) {
        // Temperature trend
        const tempDiff = tempHistory[i-1] - tempHistory[i];
        if (tempDiff > DERIVATIVE_THRESHOLD) {
            tempIncreasing++;
            tempDecreasing = 0;
        } else if (tempDiff < -DERIVATIVE_THRESHOLD) {
            tempDecreasing++;
            tempIncreasing = 0;
        }
        
        // Humidity trend
        const humidityDiff = humidityHistory[i-1] - humidityHistory[i];
        if (humidityDiff > DERIVATIVE_THRESHOLD) {
            humidityIncreasing++;
            humidityDecreasing = 0;
        } else if (humidityDiff < -DERIVATIVE_THRESHOLD) {
            humidityDecreasing++;
            humidityIncreasing = 0;
        }
    }

    return {
        increasing: tempIncreasing >= MIN_INCREASING_SAMPLES || humidityIncreasing >= MIN_INCREASING_SAMPLES,
        decreasing: tempDecreasing >= MIN_DECREASING_SAMPLES || humidityDecreasing >= MIN_DECREASING_SAMPLES,
        counts: {
            tempIncreasing,
            tempDecreasing,
            humidityIncreasing,
            humidityDecreasing
        }
    };
}

function checkThresholds(temperature, humidity, envTemp) {
    const tempDiff = temperature - envTemp;
    return {
        tempMet: tempDiff >= TEMP_THRESHOLD,
        humidityMet: humidity >= HUMIDITY_THRESHOLD
    };
}

function determineState(trends, thresholdsMet, previousState, current) {
    if(current.temperature - MINIMUM_TEMPERATURE_THRESHOLD < current.environmentalTemperature || current.humidity - MINIMUM_HUMIDITY_THRESHOLD < current.environmentalHumidity){
        return false;
    }

    // Clear decreasing trend - exit mouth state
    if (trends.decreasing) {
        return false;
    }

    // Clear increasing trend - enter mouth state
    if (trends.increasing) {
        return true;
    }
    
    // If thresholds are met, enter/maintain mouth state
    if (thresholdsMet.tempMet && thresholdsMet.humidityMet) {
        return true;
    }
    
    // If one threshold is not met, maintain previous state
    if ((thresholdsMet.tempMet || thresholdsMet.humidityMet) && previousState) {
        return true;
    }
    
    // Default to not in mouth
    return false;
}