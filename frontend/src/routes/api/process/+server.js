import { TERM_PROGRAM_VERSION } from '$env/static/private';

export async function POST({ request }) {
    try {
      const body = await request.json();
      
      const response = processData(body);
      
      return new Response(JSON.stringify({ success: true, body: response }), {
        status: 200,
        headers: { 'Content-Type': 'application/json' },
      });
    } catch (error) {
      console.error('Error handling POST request:', error);
      
      return new Response(JSON.stringify({ success: false, error: 'Invalid request' }), {
        status: 400,
        headers: { 'Content-Type': 'application/json' }
      });
    }
  }
  
  function processData(rawData) {
    const data = rawData.rawData[0];
    console.log(data);
    
    const temperature = data.humidity?.temperature || "N/A";
    const humidity = data.humidity?.humidity || "N/A";
    
    const { x, y, z } = data.accelerometer || { x: "N/A", y: "N/A", z: "N/A" };
  
    console.log(`Temperature: ${temperature} °C`);
    console.log(`Humidity: ${humidity} %`);
    console.log(`Accelerometer - X: ${x}, Y: ${y}, Z: ${z}`);

    let inMouth = false;

    //Model and variables
    const tempThresholdMouth = 27;  // Temperature in °C for the mouth
    const humidityThresholdMouth = 80; // Humidity in % for the mouth
    const tempThresholdRoom = 20; // Room temperature in °C
    const humidityThresholdRoom = 40; // Room humidity in %

    const accelerometerThreshold = 0.2; // Maximum allowed movement for mouth (in m/s²)
    const accelerometerGravity = 9.81;

    const isTemperatureInRange = temperature >= tempThresholdMouth - 2 && temperature <= tempThresholdMouth + 2;
    const isHumidityInRange = humidity >= humidityThresholdMouth - 10 && humidity <= humidityThresholdMouth + 10;
      
    const accelerometerMovement = Math.sqrt(x * x + y * y + z * z) - accelerometerGravity;
        
    const isMotionStill = Math.abs(accelerometerMovement) < accelerometerThreshold;
      
    if (isTemperatureInRange && isHumidityInRange && isMotionStill) {
        inMouth = true;
    } else {
        inMouth = false;
    }
  
    return { is_in_mouth: inMouth };
}
  