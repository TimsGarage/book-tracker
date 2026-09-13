<script lang="ts">
  import { getContext, onMount } from "svelte";
  import { Search, BookOpen, CirclePlus } from "lucide-svelte";
  import Loader from "../../../components/Loader.svelte";
  import type { Book, BookReadingStatus } from "$lib/types";
  import { handleBack } from "$lib/util";
  import type { NavState } from "../../../lib/nav_helper";
  import BookPage from "../../../components/BookPreview.svelte";
  import Input from "../../../components/Input.svelte";
  import BookCard from "../../../components/BookCard.svelte";

  let selectedBook: Book | null = $state(null);

  function resetNav() {
    navState.heading = "Search Book";
    navState.showBackButton = false;
    navState.onBack = handleBack;
    selectedBook = null;
  }

  let navState = getContext<NavState>("navState");
  if (navState) {
    resetNav();
  }

  let books = $state<Book[]>([]);
  let isLoading = $state(false);
  let errorMessage = $state("");
  let searchQuery = $state("");

  function selectBook(book: Book) {
    if (navState) {
      navState.heading = "Book Details";
      navState.showBackButton = true;
      navState.onBack = () => {
        resetNav();
      };
    }
    selectedBook = book;
  }

  function onDelete(book: Book) {
    // removeBookById(book.id);
    resetNav();
  }

  function onUpdateSave(book: Book) {
    // TODO add the update fetch here with the new book
    console.log(book);
    resetNav();
  }
</script>

<div class="page">
  {#if selectedBook != null}
    <div class="book-page">
      <BookPage
        editMode
        book={selectedBook}
        deleteCallback={onDelete}
        saveCallback={onUpdateSave}
      ></BookPage>
    </div>
  {/if}

  <div class="header">
    <Input placeholder="Search for a Book" bind:value={searchQuery}>
      {#snippet icon()}
        <Search size="24" color="var(--text)" />
      {/snippet}
    </Input>
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
        <button onclick={() => []}>Retry</button>
      </div>
    {:else if books.length === 0}
      <div class="status-view empty">
        <BookOpen size="48" color="var(--accent-color)" />
        {#if searchQuery}
          <p>No books found matching "{searchQuery}"</p>
        {:else}
          <h3>Search for a book</h3>
          <p>Add books to your library that you can't scan right now</p>
        {/if}
      </div>
    {:else}
      {#each books as book (book.id)}
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
    background-color: var(--main-background);
    z-index: 5;
  }

  .header {
    display: flex;
    flex-direction: column;
    gap: 1rem;
    padding: 2rem 1rem;
    padding-bottom: 1rem;
    /* border-bottom: 1px solid var(--outline); */
    background-color: var(--main-background);
    position: sticky;
    top: 0;
    z-index: 2;
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
