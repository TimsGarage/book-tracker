<script lang="ts">
  import {
    scan,
    cancel,
    Format,
    requestPermissions,
  } from "@tauri-apps/plugin-barcode-scanner";
  import { getContext, onDestroy, onMount } from "svelte";
  import Select from "../../../components/Select.svelte";
  import {
    Book,
    BookCheck,
    BookDown,
    BookHeart,
    Bookmark,
    BookMarked,
    BookOpenText,
    BookUp,
    Library,
    ScrollText,
  } from "lucide-svelte";
  import Datepicker from "../../../components/Datepicker.svelte";
  import type { TopNavState } from "../../../components/types";
  import { handleBack } from "../../../utils/util";
  import { lookupIsbn, createBook } from "$lib/api";
  import { goto } from "$app/navigation";
  import type { LookupBook } from "$lib/types";

  let topNavState = getContext<TopNavState>("topNavState");
  topNavState.heading = "Scanner";
  topNavState.showBackButton = true;

  let isbn = $state("");
  let isScanning = $state(false); // Track scanning state
  let isAdding = $state(false);
  let bookPromise = $state<Promise<LookupBook> | null>(null);

  async function handleAddBook(b: LookupBook) {
    isAdding = true;
    try {
      await createBook(b);
      goto("/");
    } catch (err) {
      console.error("Failed to create book:", err);
      alert("Failed to add book to library");
    } finally {
      isAdding = false;
    }
  }

  async function scanBook() {
    console.log("starting scan");
    isbn = ""; // Clear the current ISBN to show the "Scanning..." UI
    isScanning = true;

    try {
      await requestPermissions();
      const result = await scan({
        windowed: true,
        formats: [Format.EAN13],
      });

      if (
        result &&
        (result.content.length == 10 || result.content.length == 13)
      ) {
        isbn = result.content;
        bookPromise = lookupIsbn(isbn);
      }
    } catch (err) {
      console.error("Scan failed or canceled:", err);
    } finally {
      isScanning = false; // Ensure state resets whether successful, failed, or canceled
    }
  }

  async function cancelScan() {
    if (isScanning) {
      try {
        await cancel();
      } catch (err) {
        console.error("Failed to cancel the scan:", err);
      } finally {
        isScanning = false;
      }
    }

    // Fallback ISBN for testing without physical camera
    bookPromise = lookupIsbn(isbn);
  }

  // Automatically cancel the scan when navigating away from this component
  onDestroy(() => {
    cancelScan();
  });

  onMount(() => {
    scanBook();
  });

  const read_options = [
    {
      title: "Ungelesen",
      value: "to_be_read",
      icon: BookMarked,
    },
    {
      title: "Am lesen",
      value: "reading",
      icon: BookOpenText,
    },
    {
      title: "Gelesen",
      value: "read",
      icon: BookCheck,
    },
  ];

  const owning_options = [
    {
      title: "Owned",
      value: "owned",
      icon: Library,
    },
    {
      title: "Wishlist",
      value: "wishlist",
      icon: Bookmark,
    },
    {
      title: "Ausgeliehen",
      value: "lent",
      icon: BookDown,
    },
  ];

  function getTodayString(): string {
    const today = new Date();
    const year = today.getFullYear();
    const month = String(today.getMonth() + 1).padStart(2, "0");
    const day = String(today.getDate()).padStart(2, "0");
    return `${year}-${month}-${day}`;
  }
</script>

<div class="page" class:scanning={isScanning}>
  <main>
    {#if isScanning}
      <div class="scanner-overlay"></div>
    {/if}
  </main>

  <div class="popup">
    {#if isbn}
      <p class="isbn">Detected ISBN: {isbn}</p>

      {#if bookPromise}
        {#await bookPromise}
          <h2 style="margin: auto">Fetching book...</h2>
        {:then lookedUpBook}
          <div class="bookpreview">
            <span class="decoration-1"></span>
            <span class="decoration-2"></span>

            {#if lookedUpBook.thumbnail_link}
              <img
                class="thumbnail"
                src={lookedUpBook.thumbnail_link}
                alt={lookedUpBook.title || "cover"}
              />
            {/if}
            <h1 class="title">{lookedUpBook.title}</h1>
            <h3 class="author">{lookedUpBook.author}</h3>
            <p class="additional-info">
              {lookedUpBook.pages ? `${lookedUpBook.pages} Pages` : ""}
            </p>

            <p
              style="width: 100%; font-size: .8rem; color: var(--text-muted); font-weight: 600; border-bottom: 1px solid var(--outline); margin-top: .75rem; padding-block: .25rem;"
            >
              Catalogue options
            </p>

            <div class="button-group">
              <Select options={owning_options} />
              <Datepicker label="Purchased" defaultValue={getTodayString()} />
            </div>
            <Select options={read_options} fill />
          </div>

          <div class="button-group" style="margin-top: auto;">
            <button on:click={scanBook}>Rescan</button>
            <button
              class="Primary"
              disabled={isAdding}
              on:click={() => handleAddBook(lookedUpBook)}
            >
              {isAdding ? "Adding..." : "Add to Library"}
            </button>
          </div>
        {:catch err}
          <div style="text-align: center; margin: auto; padding: 1rem;">
            <p style="color: var(--text-muted); margin-bottom: 1rem;">
              {err.message || "Failed to find book details"}
            </p>
            <button on:click={scanBook} class="Primary">Try Again</button>
          </div>
        {/await}
      {/if}
    {:else if isScanning}
      <h2 style="margin: auto">Scanning...</h2>
      <button
        on:click={() => {
          cancelScan();
          handleBack();
        }}
        style="width: 100%;">Cancel</button
      >
    {:else}
      <button on:click={scanBook} style="width: 100%;" class="Primary"
        >Start Scan</button
      >
    {/if}
  </div>
</div>

<style>
  /* CRITICAL: The webview must be transparent to see the camera feed */
  :global(body, html) {
    background-color: transparent !important;
  }

  .page {
    position: relative;
    display: grid;
    height: 100%;
    width: 100%;
    grid-template-columns: auto;
    grid-template-rows: 0% 100%;
    overflow: hidden;
    transition: all 0.25s ease-out;
  }

  .page.scanning {
    grid-template-rows: 75% 25%;
  }

  main {
    position: relative;
    text-align: center;
  }

  .scanner-overlay {
    position: absolute;
    height: 120px;
    width: 220px;
    background-color: rgba(233, 233, 233, 0.1);
    transform: translate(-50%, -80%);
    top: 50%;
    left: 50%;
  }

  .page.scanning .popup {
    border-radius: 40px 40px 0px 0px;
  }

  .popup {
    position: sticky;
    height: 100%;
    bottom: 0;
    border-radius: 0px 0px;
    transition: border-radius 0.25s ease-in-out;
    background-color: var(--main-background);
    padding: 1.5rem;
    padding-bottom: 3rem;
    display: grid;
    grid-template-rows: min-content auto min-content;
  }

  .bookpreview {
    z-index: 5;
    position: relative;
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 0.5rem;
  }

  .decoration-1,
  .decoration-2 {
    z-index: 0;
    content: "";
    position: absolute;
    width: 500px;
    height: 500px;
    border-radius: 40px;
    background-color: var(--accent-color);
    rotate: 35deg;
  }

  .decoration-1 {
    left: -470px;
    top: -80px;
  }

  .decoration-2 {
    right: -480px;
    bottom: -90px;
  }

  .isbn {
    padding: 0.75rem 1.25rem;
    background-color: var(--accent-color);
    color: var(--text-reversed);
    font-weight: 600;
    border-radius: 2rem;
    font-size: 0.9rem;
    width: fit-content;
    margin-inline: auto;
  }

  .thumbnail {
    margin-block: 2rem;
    height: 270px;
  }

  .title {
    text-align: center;
    margin-inline: 1rem;
    font-size: 2.5rem;
  }

  .author,
  .additional-info {
    text-align: center;
    margin-inline: 1rem;
    opacity: 0.75;
  }

  .popup .button-group {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 0.5rem;
    width: 100%;
  }
</style>
