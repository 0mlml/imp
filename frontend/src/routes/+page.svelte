<script>
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
  
	let data = null;
	let processedData = 0;
	let pollingInterval;
  
	onMount(() => {
	  async function fetchData() {
		try {
		  const response = await fetch('http://localhost:8080/getlatest');
		  if (response.ok) {
			const result = await response.json();
			data = await processApiData(result); 
		  }
		} catch (error) {
		  console.error('Error fetching data:', error);
			processedData = -1;
		}
	  }
  
	  fetchData();
	  pollingInterval = setInterval(fetchData, 1000);
  
	  return () => clearInterval(pollingInterval);
	});
  
	async function processApiData(rawData) {
		try{
			const response = await fetch('/api/process',
				{
					method: 'POST',
					headers: {'Content-Type': 'application/json'},
					body: JSON.stringify({rawData})
				}
			);

			if (response.ok){
				processedData = (await response.json()).body.is_in_mouth;
			} 
		} catch (error) {
			console.error('Error processing data:', error);
		}
	};
  </script>
  
  <div class="min-h-screen flex flex-col items-center justify-center bg-gray-50 text-gray-800">
	<h1 class="text-4xl font-bold mb-6 text-blue-600">IMP Probe</h1>
  
	{#if processedData === 0 || processedData === -1}
	  <div class="flex items-center space-x-2">
		<div class="w-6 h-6 border-4 border-blue-500 border-t-transparent rounded-full animate-spin"></div>
		{#if processedData === 0}
			<p class="text-lg font-medium">Loading...</p>
		{:else if processedData === -1}
			<p class="text-lg font-medium">Device not responding</p>
		{/if}
	  </div>
	{:else}
	  {#if processedData === true}
		<p class="text-lg font-medium bg-green-100 text-green-700 p-4 rounded shadow">
			Device is in the mouth
		</p>
	  {:else if processedData === false}
		<p class="text-lg font-medium bg-red-100 text-red-700 p-4 rounded shadow">
			Device is not in the mouth
		</p>
	  {/if}
	{/if}

	
	<footer class="pt-10">
		<p>Developed proudly by <span class="font-bold">Team G </span></p>
	</footer>
  </div>