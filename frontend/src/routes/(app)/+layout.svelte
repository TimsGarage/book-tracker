<script lang="ts">
  import TopNav from "../../components/TopNav.svelte";
  import { setContext } from "svelte";
  import type { TopNavState } from "../../components/types";
  import BottomNav from "../../components/BottomNav.svelte";

  let { children } = $props();

  const topNavState: TopNavState = $state({
    heading: "My Library",
    showBackButton: false,
  });

  setContext("topNavState", topNavState);
</script>

<div class="page">
  <TopNav></TopNav>

  {@render children()}

  <BottomNav></BottomNav>
</div>

<style>
  .page {
    height: 100%;
    width: 100%;
    position: relative;
    display: flex;
    flex-direction: column;
  }
</style>
