<script lang="ts">
  import {
    scan,
    cancel,
    Format,
    requestPermissions,
  } from "@tauri-apps/plugin-barcode-scanner";
  import { getContext, onDestroy, onMount } from "svelte";
  import type { NavState } from "../../../lib/nav_helper";
  import { handleBack, owning_options, read_options } from "../../../lib/util";
  import { lookupIsbn, createBook } from "$lib/api";
  import { goto } from "$app/navigation";
  import type { LookupBook, Book } from "$lib/types";
  import BookPage from "../../../components/BookPreview.svelte";
  import Loader from "../../../components/Loader.svelte";
  import { MonitorX } from "lucide-svelte";

  let navState = getContext<NavState>("navState");
  navState.heading = "Scanner";
  navState.showBackButton = true;
  navState.hideBottomNavigation = true;

  let isbn = $state("");
  let isScanning = $state(false); // Track scanning state
  let isAdding = $state(false);
  let bookPromise = $state<Promise<LookupBook> | null>(null);

  function getFixed() {
    isbn = "978-3-426-65443-9";
    bookPromise = lookupIsbn(isbn);
  }

  async function handleAddBook(b: LookupBook) {
    isAdding = true;
    try {
      await createBook(b);
      goto("/");
    } catch (err: any) {
      console.error("Failed to create book:", err);
      alert(err.message || "Failed to add book to library");
    } finally {
      isAdding = false;
    }
  }

  async function scanBook() {
    console.log("starting scan");
    isbn = "";
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
      isScanning = false;
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

    bookPromise = lookupIsbn(isbn);
  }

  onDestroy(() => {
    cancelScan();
  });

  onMount(() => {
    scanBook();
  });
</script>

<div class="page" class:scanning={isScanning}>
  <main>
    {#if isScanning}
      <div class="scanner-overlay"></div>
    {/if}
  </main>

  <div class="popup">
    {#if isbn}
      {#if bookPromise}
        {#await bookPromise}
          <div class="status-view">
            <Loader size="48px"></Loader>
            <p>Looking for your book...</p>
          </div>
        {:then lookedUpBook}
          <BookPage
            showSaveButton
            book={lookedUpBook}
            saveCallback={(book: Book) => {
              handleAddBook(book);
            }}
            deleteCallback={scanBook}
          ></BookPage>
        {:catch err}
          <div style="text-align: center; margin: auto; padding: 1rem;">
            <p style="color: var(--text-muted); margin-bottom: 1rem;">
              {err.message || "Failed to find book details"}
            </p>
            <button onclick={scanBook} class="Primary">Try Again</button>
          </div>
        {/await}
      {/if}
    {:else if isScanning}
      <button
        onclick={() => {
          cancelScan();
          handleBack();
        }}
        style="margin: 2rem; margin-bottom: 0px;">Cancel</button
      >
    {:else}
      <div class="status-view">
        <MonitorX color="var(--text)" size="48"></MonitorX>
        <h5>Something went wrong.</h5>
        <button
          onclick={scanBook}
          style="width: 100%; min-height: unset; padding: .75rem 1.5rem;"
          class="Primary">Try again</button
        >
      </div>
      <!-- TODO dev stuff -->
      <!-- <button onclick={getFixed} style="width: 100%;">Get Fixed</button> -->
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
    grid-template-rows: auto min-content;
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
    padding-bottom: 2rem;
    display: grid;
    grid-template-rows: calc(100% - 64px) 64px;
  }

  .status-view {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    margin: auto;
    padding: 4rem;
    gap: 1rem;
    text-align: center;
    color: var(--text-muted);
  }
</style>
