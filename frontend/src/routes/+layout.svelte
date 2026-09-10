<script>
  let { children } = $props();
</script>

<div class="layout">
  {@render children()}
</div>

<style>
  .layout {
    position: relative;
    display: flex;
    flex-direction: column;
    height: 100dvh;
    width: 100dvw;
  }
</style>
