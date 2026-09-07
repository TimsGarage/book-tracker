<script>
  import { scan, Format, requestPermissions } from '@tauri-apps/plugin-barcode-scanner';
  let isbn = "";

  async function scanBook() {
    try {
      // Request Android camera permissions first
      await requestPermissions();
      
      // Filter specifically for EAN-13, which is the standard format for book ISBNs
      const result = await scan({ 
        windowed: true, 
        formats: [Format.EAN13] 
      });
      
      isbn = result.content;
    } catch (err) {
      console.error("Scan failed or canceled:", err);
    }
  }
</script>

<main>
  <h1>ISBN Scanner</h1>
  <button on:click={scanBook}>Open Camera</button>
  
  {#if isbn}
    <div class="result">
      <p><strong>Scanned ISBN:</strong> {isbn}</p>
    </div>
  {/if}
</main>

<style>
  /* CRITICAL: The webview must be transparent to see the camera feed */
  :global(body, html) {
    background-color: transparent !important;
  }
  
  main {
    padding: 2rem;
    text-align: center;
    background: rgba(255, 255, 255, 0.8); /* Semi-transparent overlay */
    border-radius: 12px;
    margin: 20px;
  }
</style>