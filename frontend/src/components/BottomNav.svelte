<script lang="ts">
  import {
    Bookmark,
    BookSearch,
    ChevronLeft,
    CircleUserRound,
    Library,
    Scan,
    ScanBarcode,
    Settings,
  } from "lucide-svelte";
</script>

<nav class="bottomnav">
  {#snippet navItem(Icon: ConstructorOfATypedSvelteComponent, link: string)}
    <a class="navItem" href={link}>
      <Icon></Icon>
    </a>
  {/snippet}
  <div class="content">
    {@render navItem(Library, "/")}
    {@render navItem(Bookmark, "/")}
    <a class="scanner" href="/scanner">
      <ScanBarcode></ScanBarcode>
    </a>
    {@render navItem(BookSearch, "/")}
    {@render navItem(Settings, "/login")}
  </div>
</nav>

<style>
  .bottomnav {
    position: sticky;
    bottom: 0;
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
