<script lang="ts">
  import { ChevronLeft, CircleUserRound, Library } from "lucide-svelte";
  import { getContext } from "svelte";
  import type { TopNavState } from "./types";
  import { handleBack } from "../utils/util";

  let topNavState = getContext<TopNavState>("topNavState");
</script>

<nav class="topnav">
  <div class="content">
    {#if topNavState.showBackButton}
      <!-- svelte-ignore a11y_click_events_have_key_events -->
      <!-- svelte-ignore a11y_no_static_element_interactions -->
      <span on:click={handleBack} style="height: 32px"
        ><ChevronLeft size="32" /></span
      >
    {:else}
      <a href="/"><Library size="32" /></a>
    {/if}
    <h3>{topNavState.heading}</h3>

    {#if topNavState.showBackButton}
      <span></span>
    {:else}
      <span class="decoration-1"></span>
      <a href="/dev"
        ><CircleUserRound size="32" color="var(--main-background)" /></a
      >
    {/if}
  </div>
</nav>

<style>
  .topnav {
    position: sticky;
    top: 0;
    width: 100%;
    background-color: var(--main-background);
  }

  a {
    color: var(--accent-color);
  }

  .content {
    position: relative;
    border-bottom: 1px solid var(--outline);
    padding-top: 3rem; /* Padding to get rid of phone notch */
    padding-inline: 1.5rem;
    display: grid;
    grid-template-columns: 32px auto 32px;
    height: calc(3rem + 70px);
    align-items: center;
    gap: 0.5rem;
    overflow: hidden;
  }

  .content * {
    z-index: 5;
  }

  .decoration-1 {
    z-index: 0;
    position: absolute;
    width: 300px;
    height: 300px;
    background-color: var(--accent-color);
    rotate: 130deg;
  }

  .decoration-1 {
    right: -190px;
    top: 30px;
  }
</style>
