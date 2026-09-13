<script lang="ts">
  import { getContext, onMount } from "svelte";
  import { Search, BookOpen, CirclePlus, Loader, Star } from "lucide-svelte";
  import { fetchWishlistBooks } from "$lib/api";
  import type { Book, BookReadingStatus } from "$lib/types";
  import { handleBack } from "$lib/util";
  import type { NavState } from "../../../lib/nav_helper";
  import BookCard from "../../../components/BookCard.svelte";
  import { goto } from "$app/navigation";

  let navState = getContext<NavState>("navState");
  if (navState) {
    navState.heading = "Wishlist";
    navState.showBackButton = false;
    navState.icon = Star;
    navState.onBack = handleBack;
    navState.hideBottomNavigation = false;
  }

  let books = $state<Book[]>([]);
  let isLoading = $state(true);
  let errorMessage = $state("");
  let searchQuery = $state("");

  async function loadBooks() {
    isLoading = true;
    errorMessage = "";
    try {
      books = await fetchWishlistBooks();
    } catch (err: any) {
      console.error("Failed to load books:", err);
      errorMessage = err.message || "Failed to load wishlist";
    } finally {
      isLoading = false;
    }
  }

  onMount(() => {
    loadBooks();
  });

  // 1. Filter books ONLY by search query
  const filteredBooks = $derived(
    books.filter((b: Book) => {
      const q = searchQuery.toLowerCase().trim();
      if (!q) return true;
      return (
        (b.title || "").toLowerCase().includes(q) ||
        (b.author || "").toLowerCase().includes(q) ||
        (b.isbn || "").toLowerCase().includes(q)
      );
    }),
  );
</script>

<div class="page">
  <!-- <div class="header">
    <Input placeholder="Search for a Book" bind:value={searchQuery}>
      {#snippet icon()}
        <Search size="24" color="var(--text)" />
      {/snippet}
    </Input>
  </div> -->

  <div class="books-container">
    {#if isLoading}
      <div class="status-view">
        <Loader size="48px" />
        <p>Loading your Wishlist...</p>
      </div>
    {:else if errorMessage}
      <div class="status-view error">
        <p>{errorMessage}</p>
        <button onclick={loadBooks}>Retry</button>
      </div>
    {:else if filteredBooks.length === 0}
      <div class="status-view empty">
        <BookOpen size="48" color="var(--accent-color)" />
        {#if searchQuery}
          <p>No books found matching "{searchQuery}"</p>
        {:else}
          <h3>Your wishlist is empty</h3>
          <p>Scan a book to start your list</p>
          <a href="/scanner" class="btn-primary">
            <CirclePlus size="18" />
            <span>Scan a Book</span>
          </a>
        {/if}
      </div>
    {:else}
      {#each filteredBooks as book (book.id)}
        <BookCard onclick={() => goto("/" + book.id)} {book} />
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

  .header {
    display: flex;
    flex-direction: column;
    gap: 1rem;
    padding: 2rem 1rem;
    padding-bottom: 1rem;
    border-bottom: 1px solid var(--outline);
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
    padding: 4rem 2rem;
    gap: 1rem;
    text-align: center;
    color: var(--text-muted);
  }

  .status-view.empty h3 {
    color: var(--text);
    margin: 0;
  }

  .btn-primary {
    display: inline-flex;
    align-items: center;
    gap: 0.5rem;
    background-color: var(--accent-color);
    color: var(--text-reversed);
    padding: 0.75rem 1.5rem;
    border-radius: 0.5rem;
    text-decoration: none;
    font-weight: 600;
    margin-top: 0.5rem;
  }
</style>
