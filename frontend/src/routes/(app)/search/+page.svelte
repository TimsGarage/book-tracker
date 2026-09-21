<script lang="ts">
  import { page } from "$app/state";
  import { goto } from "$app/navigation";
  import { getContext, onDestroy } from "svelte";
  import { Search, BookOpen } from "lucide-svelte";
  import Loader from "../../../components/Loader.svelte";
  import type { Book, LookupBook } from "$lib/types";
  import type { NavState } from "../../../lib/nav_helper";
  import BookPage from "../../../components/BookPreview.svelte";
  import Input from "../../../components/Input.svelte";
  import BookCard from "../../../components/BookCard.svelte";
  import { createBook, lookupSearchterm } from "$lib/api";

  let navState = getContext<NavState>("navState");

  let books = $state<LookupBook[]>([]);
  let isLoading = $state(false);
  let errorMessage = $state("");
  let searchQuery = $state("");

  let selectedIsbn = $derived(page.url.searchParams.get("isbn"));

  let selectedBook = $derived(
    selectedIsbn ? (books.find((b) => b.isbn === selectedIsbn) ?? null) : null,
  );

  $effect(() => {
    if (!navState) return;

    if (selectedBook) {
      navState.hideBottomNavigation = true;
      navState.heading = "Book Details";
      navState.showBackButton = true;
      navState.onBack = () => history.back();
    } else {
      navState.hideBottomNavigation = false;
      navState.heading = "Search Book";
      navState.showBackButton = false;
      navState.onBack = undefined;
    }
  });

  function selectBook(book: Book | LookupBook) {
    const nextUrl = new URL(page.url);
    nextUrl.searchParams.set("isbn", book.isbn);
    goto(nextUrl.toString(), { keepFocus: true, noScroll: true });
  }

  async function onUpdateSave(book: Book) {
    try {
      goto("/");
      await createBook(book);
    } catch (err: any) {
      console.error("Failed to create book:", err);
      alert(err.message || "Failed to add book to library");
    }
  }

  async function search() {
    isLoading = true;
    errorMessage = "";
    try {
      let res = await lookupSearchterm(searchQuery);
      books = res.filter((book) => book.isbn !== "");
    } catch (err) {
      errorMessage = "Search failed. Please try again.";
    } finally {
      isLoading = false;
    }
  }

  onDestroy(() => {
    if (navState) {
      navState.heading = "Search Book";
      navState.showBackButton = false;
      navState.onBack = undefined;
    }
  });
</script>

<div class="page">
  {#if selectedBook}
    <div class="book-page">
      <BookPage
        showSaveButton
        book={selectedBook}
        saveCallback={onUpdateSave}
      />
    </div>
  {/if}

  <div class="header">
    <Input
      placeholder="Search for a Book"
      bind:value={searchQuery}
      onkeypress={(e) => {
        if (e.key == "Enter") {
          search();
        }
      }}
    ></Input>
    <button class="serach-button Primary" onclick={search}
      ><Search></Search></button
    >
  </div>

  <div class="books-container">
    {#if isLoading}
      <div class="status-view">
        <Loader size="48px" />
        <p>Searching...</p>
      </div>
    {:else if errorMessage}
      <div class="status-view error">
        <p>{errorMessage}</p>
        <button onclick={search}>Retry</button>
      </div>
    {:else if books.length === 0}
      <div class="status-view empty">
        <BookOpen size="48" color="var(--accent-color)" />
        <h3>Search for a book</h3>
        <p>Add books to your library that you can't scan right now</p>
      </div>
    {:else}
      {#each books as book (book)}
        <BookCard onclick={() => selectBook(book)} {book} />
      {/each}
    {/if}
  </div>
</div>

<style>
  .page {
    height: 100%;
    width: 100%;
    position: relative;
    display: flex;
    flex-direction: column;
    background-color: var(--main-background);
    overflow-y: auto;
  }

  .book-page {
    position: absolute;
    top: 0;
    left: 0;
    height: 100%;
    width: 100%;
    background-color: var(--main-background);
    z-index: 5;
  }

  .header {
    display: grid;
    grid-template-columns: auto min-content;
    gap: 0.25rem;
    padding: 2rem 1rem;
    padding-bottom: 1rem;
    background-color: var(--main-background);
    position: sticky;
    top: 0;
    z-index: 2;
  }

  .serach-button {
    min-height: initial;
    min-width: initial;
    aspect-ratio: 1/1;
  }

  .books-container {
    flex: 1;
    overflow-y: auto;
  }

  .status-view {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 4rem;
    gap: 1rem;
    text-align: center;
    color: var(--text-muted);
    height: 90%;
  }

  .status-view.empty h3 {
    color: var(--text);
    margin: 0;
  }
</style>
