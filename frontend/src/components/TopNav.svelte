<script lang="ts">
  import { ChevronLeft, CircleUserRound, Library } from "lucide-svelte";
  import { getContext } from "svelte";
  import type { NavState } from "../lib/nav_helper";
  import { handleBack } from "../lib/util";

  let navState = getContext<NavState>("navState");

  let Icon = $derived(navState.icon ?? Library);
</script>

<nav class="topnav">
  <div class="content">
    {#if navState.showBackButton}
      <!-- svelte-ignore a11y_click_events_have_key_events -->
      <!-- svelte-ignore a11y_no_static_element_interactions -->
      <span
        onclick={navState.onBack ? navState.onBack() : handleBack()}
        style="height: 32px"><ChevronLeft size="32" /></span
      >
    {:else if navState.icon}
      <a href="/" class="icon"><Icon size="28" /></a>
    {/if}
    <h3>{navState.heading}</h3>
  </div>
</nav>

<style>
  .topnav {
    position: sticky;
    top: 0;
    width: 100%;
    background-color: var(--main-background);
  }

  .icon {
    height: 28px;
    position: relative;
    color: var(--main-background);
  }

  .icon::before {
    z-index: -1;
    content: "";
    position: absolute;
    width: calc(1.5rem + 30px + 0.75rem);
    height: calc(30px + 1.25rem);
    background-color: var(--accent-color);
    left: -1.5rem;
    top: -0.625rem;
    border-radius: 0px 1.5rem 1.5rem 0px;
  }

  .content {
    position: relative;
    border-bottom: 1px solid var(--outline);
    padding-top: 3rem; /* Padding to get rid of phone notch */
    padding-inline: 1.5rem;
    display: flex;
    height: calc(3rem + 70px);
    align-items: center;
    gap: 1.25rem;
    overflow: hidden;
  }

  .content * {
    z-index: 5;
  }
</style>
