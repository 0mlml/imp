export async function GET() {
    await new Promise(r => setTimeout(r, 1000));

    const data = await (await fetch(`http://localhost:8080/getlatest?count=50`)).json();

    const responseBody = getEnvironmentalData(data);

    return new Response(JSON.stringify({ success: true, body: responseBody}), {
        status: 200,
        headers: { 'Content-Type': 'application/json' },
    });
}

const getEnvironmentalData = function (data) {
    let temperature = 0;
    let humidity = 0;
    let count = 0;
    let i = 0;
    while (data[i] != undefined) {
        temperature += data[i].humidity.temperature;
        humidity += data[i].humidity.humidity;
        count += 1;
        i++;
    }

    temperature = temperature / count;
    humidity = humidity / count
    return { environmentalTemperature: temperature, environmentalHumidity: humidity };
}