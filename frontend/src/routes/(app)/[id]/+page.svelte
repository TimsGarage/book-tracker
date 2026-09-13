<script lang="ts">
  import { fetchBookById, removeBookById, updateBook } from "$lib/api";
  import type { PageProps } from "./$types";
  let { params }: PageProps = $props();
  import { getContext, onMount } from "svelte";
  import { type NavState } from "../../../lib/nav_helper";
  import Loader from "../../../components/Loader.svelte";
  import type { Book } from "$lib/types";
  import BookPage from "../../../components/BookPreview.svelte";
  import { handleBack, owning_options } from "$lib/util";
  let navState = getContext<NavState>("navState");
  if (navState) {
    navState.heading = "Book Details";
    navState.showBackButton = true;
    navState.onBack = handleBack;
    navState.hideBottomNavigation = true;
  }

  let book = $state<Book>();
  let isLoading = $state(true);
  let errorMessage = $state("");

  async function loadBook() {
    isLoading = true;
    errorMessage = "";
    try {
      book = await fetchBookById(params.id as unknown as number);
    } catch (err: any) {
      console.error("Failed to load book:", err);
      errorMessage = err.message || "Failed to load library";
    } finally {
      isLoading = false;
    }
  }

  onMount(() => {
    loadBook();
  });

  function onDelete(book: Book) {
    removeBookById(book.id);
    handleBack();
  }

  function onUpdateSave(book: Book) {
    updateBook(book);
    handleBack();
  }
</script>

<div class="page">
  {#if isLoading}
    <div class="status-view">
      <Loader size="48px" />
      <p>Loading your Book...</p>
    </div>
  {:else if errorMessage}
    <div class="status-view error">
      <p>{errorMessage}</p>
      <button onclick={loadBook}>Retry</button>
    </div>
  {:else if book != null}
    <BookPage
      showDeleteButton
      showSaveButton
      {book}
      deleteCallback={onDelete}
      saveCallback={onUpdateSave}
    ></BookPage>
  {/if}
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

  .status-view {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 4rem;
    height: 90%;
    gap: 1rem;
    text-align: center;
    color: var(--text-muted);
  }
</style>
