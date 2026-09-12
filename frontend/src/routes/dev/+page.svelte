<script lang="ts">
  import Loader from "../../components/Loader.svelte";
  import BottomNav from "../../components/BottomNav.svelte";
</script>

<div class="page">
  <h1>Development page</h1>
  <div class="content">
    <div class="card components-card">
      <h3>Design System Preview</h3>
      <div class="loader-demo">
        <Loader size="48px" />
      </div>
    </div>
  </div>

  <BottomNav />
</div>

<style>
  .page {
    position: relative;
    display: flex;
    flex-direction: column;
    height: 100%;
    background-color: var(--main-background);
  }

  .card {
    background-color: var(--secondary-background);
    border: 1px solid var(--outline);
    border-radius: 1rem;
    padding: 1.5rem;
    margin: 1.5rem;
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }
  .components-card h3 {
    margin: 0;
    font-size: 1rem;
    color: var(--text-muted);
  }

  .loader-demo {
    display: flex;
    justify-content: center;
    padding: 1rem;
  }
</style>
