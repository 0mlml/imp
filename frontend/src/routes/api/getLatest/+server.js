import { GRAVITY_ACCEL, HUM_MAX_DELTA, TEMP_DIFFERENCE, TEMP_MAX_DELTA, THRESHOLD_ACCEL, THRESHOLD_HUMIDITY_DIFFERENCE, THRESHOLD_HUMIDITY_IN_MOUTH, THRESHOLD_STILL_PERIOD, THRESHOLD_STRICT_ACCEL, THRESHOLD_TEMP_DIFFERENCE, THRESHOLD_TEMP_IN_MOUTH } from "$lib/consts";

export async function POST({ request }) {
    const { environmentalTemperature, environmentalHumidity } = await request.json();
    const data = await (await fetch(`http://localhost:8080/getlatest?count=100`)).json();

    const { humidity, temperature } = data[0].humidity;
    const { x, y, z, peak_acceleration } = data[0].accelerometer;

    const isTemperatureInRange = temperatureInMouth(temperature, environmentalTemperature);
    const isHumidityInRange = humidityInMouth(humidity, environmentalHumidity);
    const isMotionStill = Math.abs(peak_acceleration - GRAVITY_ACCEL) <= THRESHOLD_ACCEL;
    const inMouth = (isTemperatureInRange ? 0.6 : -0.2) + (isHumidityInRange ? 0.6 : -0.3) + (isMotionStill ? -0.5 : 0.3) > 1; 
    return new Response(JSON.stringify({ success: true, body: {environmentalHumidity, environmentalTemperature, temperature, humidity, x, y, z, inMouth, isTemperatureInRange, isTemperatureInRange, isMotionStill , pa: (peak_acceleration-GRAVITY_ACCEL)} }), {
        status: 200,
        headers: { 'Content-Type': 'application/json' },
    });
}

function temperatureInMouth(temperature, environmentalTemperature) {
    const diff = temperature - environmentalTemperature;
    const aboveThreshold = temperature >= THRESHOLD_TEMP_IN_MOUTH;
    return (aboveThreshold ? 1 : -0.5) + (Math.abs(diff) > THRESHOLD_TEMP_DIFFERENCE ? 0.5 : 0) + (diff < 0 ? -10 : 0.2) > 1;
}

function humidityInMouth(humidity, environmentalHumidity) {
    const diff = humidity - environmentalHumidity;
    const aboveThreshold = humidity >= THRESHOLD_HUMIDITY_IN_MOUTH;
    return (aboveThreshold ? 5 : -0.5) + (Math.abs(diff) > THRESHOLD_HUMIDITY_DIFFERENCE ? 0.5 : 0) + (diff < 0 ? -0.5 : 0.5) > 1;
}

function calculateEnvironmentalReadings(period) {
    const temps = period.map(r => r.humidity.temperature);
    const hums = period.map(r => r.humidity.humidity);

    return {
        environmentalTemperature: median(temps),
        environmentalHumidity: median(hums)
    };
}

function median(values) {
    const sorted = [...values].sort((a, b) => a - b);
    const mid = Math.floor(sorted.length / 2);
    return sorted.length % 2 ? sorted[mid] : (sorted[mid - 1] + sorted[mid]) / 2;
}

function smoothValue(last, current) {
    return last * 0.8 + current * 0.2;
}