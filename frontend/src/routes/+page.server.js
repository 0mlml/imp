export async function load(event) {
    try {
        const response = await event.fetch('/api/getEnvironment');
        if (response.ok) {
        const environmentData = (await response.json()).body;
        return {environmentData: environmentData};
        }
      } catch (error) {
        console.error('Error updating environment data:', error);
      }
}