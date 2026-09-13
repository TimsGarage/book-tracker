<script lang="ts">
  import TopNav from "../../components/TopNav.svelte";
  import { setContext } from "svelte";
  import BottomNav from "../../components/BottomNav.svelte";
  import { handleBack } from "$lib/util";
  import type { NavState } from "$lib/nav_helper";

  let { children } = $props();

  const navState: NavState = $state({
    heading: "My Library",
    showBackButton: false,
    onBack: () => handleBack,
  });

  setContext("navState", navState);
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
