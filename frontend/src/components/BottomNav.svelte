<script lang="ts">
  import {
    Bookmark,
    BookSearch,
    ChevronLeft,
    CircleUserRound,
    Code,
    Library,
    Scan,
    ScanBarcode,
    Settings,
    Star,
  } from "lucide-svelte";
  import BottomNavItem from "./BottomNavItem.svelte";
  import { getContext } from "svelte";
  import type { NavState } from "$lib/nav_helper";

  let navState = getContext<NavState>("navState");
</script>

{#if !navState.hideBottomNavigation}
  <nav class="bottomnav">
    <div class="content">
      <BottomNavItem Icon={Library} link="/"></BottomNavItem>
      <BottomNavItem Icon={Star} link="/wishlist"></BottomNavItem>
      <a class="scanner" href="/scanner">
        <ScanBarcode></ScanBarcode>
      </a>
      <BottomNavItem Icon={BookSearch} link="/search"></BottomNavItem>
      <BottomNavItem Icon={CircleUserRound} link="/profile"></BottomNavItem>
    </div>
  </nav>
{/if}

<style>
  .bottomnav {
    z-index: 10;
    position: sticky;
    top: 100%;
    bottom: 0;
    left: 0;
    width: 100%;
    background-color: var(--main-background);
  }

  .content {
    position: relative;
    border-top: 1px solid var(--outline);
    /* padding-top: 3rem; Padding to get rid of phone notch */
    padding: 1rem;
    padding-top: 0;
    height: calc(80px);
    overflow: show;
    display: flex;
    justify-content: space-around;
    align-items: center;
  }

  .scanner,
  .navItem {
    width: 52px;
    height: 52px;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .navItem.active {
    background-color: red;
  }

  .scanner {
    position: relative;
    top: -30%;
    background-color: var(--accent-color);
    color: var(--main-background);
    border-radius: 100%;
    width: 72px;
    height: 72px;
  }
</style>
