<script>
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { DATA_STATE_LOADING, DATA_STATE_DEVICE_NOT_RESPONDING, DATA_STATE_IN_MOUTH, DATA_STATE_OUT_MOUTH} from '$lib/consts.js';
  import { env } from '$env/dynamic/public';

  let dataState = $state(DATA_STATE_LOADING);
	let processedData = $state(null);
	let pollingInterval;

  let updateEnvironment = true;
  
	onMount(() => {
	  async function fetchData() {
      try {
        const response = await fetch('/api/getLatest', {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify({updateEnvironment: updateEnvironment})
        });
        if (response.ok) {
        processedData = (await response.json()).body;
        dataState = processedData.inMouth ? DATA_STATE_IN_MOUTH : DATA_STATE_OUT_MOUTH;
        updateEnvironment = false;
        }
      } catch (error) {
        console.error('Error fetching data:', error);
        dataState = DATA_STATE_DEVICE_NOT_RESPONDING;
      }
	  }
  
    fetchData();
	  pollingInterval = setInterval(fetchData, 300);
  
	  return () => clearInterval(pollingInterval);
	});

  function onclick(){
    updateEnvironment = true;
  }
   
  </script>
  
  <div class="min-h-screen flex flex-col items-center justify-center bg-gray-50 text-gray-800">
	<h1 class="text-4xl font-bold mb-6 text-blue-600">IMP Probe</h1>
  
	{#if dataState === DATA_STATE_DEVICE_NOT_RESPONDING || dataState === DATA_STATE_LOADING}
	  <div class="flex items-center space-x-2">
        <div class="w-6 h-6 border-4 border-blue-500 border-t-transparent rounded-full animate-spin"></div>
            {#if dataState === DATA_STATE_LOADING}
                <p class="text-lg font-medium">Loading...</p>
            {:else if dataState === DATA_STATE_DEVICE_NOT_RESPONDING}
                <p class="text-lg font-medium">Device not responding</p>
            {/if}
	  </div>
  {:else}
  {#if dataState === DATA_STATE_IN_MOUTH}
  <p class="text-lg font-medium bg-green-100 text-green-700 p-4 rounded mb-6 shadow hover:shadow-lg">
    Device is in the mouth
  </p>
{:else if dataState === DATA_STATE_OUT_MOUTH}
  <p class="text-lg font-medium bg-red-100 text-red-700 p-4 rounded mb-6 shadow hover:shadow-lg">
    Device is not in the mouth
  </p>
{/if}

<p class="text-center mb-4">
  <button
    class="border border-blue-200 bg-blue-50 text-blue-600 px-4 py-2 rounded shadow-md hover:bg-blue-100 hover:scale-110 hover:shadow-lg transition-transform duration-200"
    {onclick}>
    Update Environment
  </button>
</p>
      
    <div class="border border-gray-300 bg-white p-4 rounded shadow-xl mb-6 w-72"> <!-- position -->
      <h2 class="font-bold text-lg mb-2 text-center">Position</h2>
      <p class="mb-2">
        <span class="font-bold italic">X:</span> {processedData.x.toFixed(2)} m/s²
      </p>
      <p class="mb-2">
        <span class="font-bold italic">Y:</span> {processedData.y.toFixed(2)} m/s²
      </p>
      <p class="mb-2">
        <span class="font-bold italic">Z:</span> {processedData.z.toFixed(2)} m/s²
      </p>
      <p>
        <span class="font-bold italic">Still:</span> {processedData.isMotionStill}
      </p>
    </div>

      <div class="border border-gray-300 bg-white p-4 rounded shadow-xl mb-6 w-72"> <!-- temp and humidity -->
  <h2 class="font-bold text-lg mb-2 text-center">Temperature & Humidity</h2>
  <p class="mb-2">
    <span class="font-bold italic">Temperature:</span> {processedData.temperature.toFixed(2)} °C
  </p>
  <p class="mb-2">
    <span class="font-bold italic">Humidity:</span> {processedData.humidity.toFixed(2)} %
  </p>
</div>

<div class="border border-gray-300 bg-white p-4 rounded shadow-xl mb-6 w-72"> <!-- environment -->
  <h2 class="font-bold text-lg mb-2 text-center">Environment</h2>
  <p class="mb-2">
    <span class="font-bold italic">eTemperature:</span> {processedData.environmentalTemperature.toFixed(2)} °C
  </p>
  <p class="mb-2">
    <span class="font-bold italic">eHumidity:</span> {processedData.environmentalHumidity.toFixed(2)} %
  </p>
  <p>
    <span class="font-bold italic">Peak Acceleration:</span> {processedData.peak_acceleration.toFixed(2)} m/s²
  </p>
</div>

	{/if}

	
	<footer class="pt-10">
		<p>Developed proudly by <span class="font-bold">Team G</span>🔥</p>
	</footer>
  </div>
