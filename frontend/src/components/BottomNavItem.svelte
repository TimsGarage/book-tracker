<script lang="ts">
  import { page } from "$app/state";
  import type { Component } from "svelte";

  const {
    Icon,
    link,
  }: {
    Icon: Component;
    link: string;
  } = $props();

  let active = $derived(page.url.pathname === link);
</script>

<a class="navItem" href={link} class:active>
  <Icon />
</a>

<style>
  .navItem {
    position: relative;
    width: 52px;
    height: 52px;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .navItem.active::after {
    content: "";
    position: absolute;
    width: 60%;
    height: 3px;
    border-radius: 5px;
    background: var(--accent-color);
    bottom: 5px;
    left: 20%;
  }

  .navItem.active {
    color: var(--accent-color);
  }
</style>
