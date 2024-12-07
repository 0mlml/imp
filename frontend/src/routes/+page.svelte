<script>
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { DATA_STATE_LOADING, DATA_STATE_DEVICE_NOT_RESPONDING, DATA_STATE_IN_MOUTH, DATA_STATE_OUT_MOUTH} from '$lib/consts.js';

    let dataState = DATA_STATE_LOADING;
	let processedData = null;
	let pollingInterval;
  
	onMount(() => {
	  async function fetchData() {
      try {
        const response = await fetch('/api/getLatest');
        if (response.ok) {
        processedData = (await response.json()).body;
        dataState = processedData.inMouth ? DATA_STATE_IN_MOUTH : DATA_STATE_OUT_MOUTH;
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
		<p class="text-lg font-medium bg-green-100 text-green-700 p-4 rounded shadow">
			Device is in the mouth
		</p>
	  {:else if dataState === DATA_STATE_OUT_MOUTH}
		<p class="text-lg font-medium bg-red-100 text-red-700 p-4 rounded shadow">
			Device is not in the mouth
		</p>
	  {/if} 
      <div><!-- position -->
        <p class="text-lg font-medium">X: {processedData.x}</p>
        <p class="text-lg font-medium">Y: {processedData.y}</p>
        <p class="text-lg font-medium">Z: {processedData.z}</p>
        <p class="text-lg font-medium">Still?: {processedData.isMotionStill}</p>
      </div>
      <div> <!-- temp humidity -->
        <p class="text-lg font-medium">Temperature: {processedData.temperature}</p>
        <p class="text-lg font-medium">Humidity: {processedData.humidity}</p>
      </div>
      <div> <!-- environment -->
        <p class="text-lg font-medium">eTemperature: {processedData.environmentalTemperature}</p>
        <p class="text-lg font-medium">eHumidity: {processedData.environmentalHumidity}</p>
        <p class="text-lg font-medium">peak accel: {processedData.pa}</p>
      </div>
	{/if}

	
	<footer class="pt-10">
		<p>Developed proudly by <span class="font-bold">Team G</span></p>
	</footer>
  </div>