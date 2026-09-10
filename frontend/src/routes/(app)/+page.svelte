<script lang="ts">
  import { getContext, onMount } from "svelte";
  import BookCard from "../../components/BookCard.svelte";
  import type { TopNavState } from "../../components/types";
  import Input from "../../components/Input.svelte";
  import Loader from "../../components/Loader.svelte";
  import { Search, PlusCircle, BookOpen } from "lucide-svelte";
  import Chip from "../../components/Chip.svelte";
  import { fetchBooks } from "$lib/api";
  import type { Book } from "$lib/types";

  let topNavState = getContext<TopNavState>("topNavState");
  if (topNavState) {
    topNavState.heading = "Home";
    topNavState.showBackButton = false;
  }

  let books = $state<Book[]>([]);
  let isLoading = $state(true);
  let errorMessage = $state("");
  let searchQuery = $state("");
  let activeFilter = $state<"all" | "unread" | "read" | "reading">("all");

  async function loadBooks() {
    isLoading = true;
    errorMessage = "";
    try {
      books = await fetchBooks();
    } catch (err: any) {
      console.error("Failed to load books:", err);
      errorMessage = err.message || "Failed to load library";
    } finally {
      isLoading = false;
    }
  }

  onMount(() => {
    loadBooks();
  });

  const filteredBooks = $derived(
    books.filter((b) => {
      const q = searchQuery.toLowerCase().trim();
      const matchesSearch =
        !q ||
        (b.Title || "").toLowerCase().includes(q) ||
        (b.Author || "").toLowerCase().includes(q) ||
        (b.Isbn || "").toLowerCase().includes(q);

      if (!matchesSearch) return false;

      if (activeFilter === "unread") return !b.read;
      if (activeFilter === "read") return !!b.read;
      return true;
    }),
  );

  const unreadCount = $derived(books.filter((b) => !b.read).length);
  const readCount = $derived(books.filter((b) => !!b.read).length);
</script>

<div class="page">
  <div class="header">
    <Input placeholder="Search for a Book" bind:value={searchQuery}>
      {#snippet icon()}
        <Search size="24" color="var(--text)" />
      {/snippet}
    </Input>

    <div class="tags">
      <button
        type="button"
        class="tag-btn"
        onclick={() => (activeFilter = "all")}
      >
        <Chip
          text="All Books ({books.length})"
          selected={activeFilter === "all"}
        />
      </button>
      <button
        type="button"
        class="tag-btn"
        onclick={() => (activeFilter = "unread")}
      >
        <Chip
          text="Ungelesen ({unreadCount})"
          selected={activeFilter === "unread"}
        />
      </button>
      <button
        type="button"
        class="tag-btn"
        onclick={() => (activeFilter = "read")}
      >
        <Chip text="Gelesen ({readCount})" selected={activeFilter === "read"} />
      </button>
    </div>
  </div>

  <div class="books-container">
    {#if isLoading}
      <div class="status-view">
        <Loader size="48px" />
        <p>Loading your library...</p>
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
          <h3>Your library is empty</h3>
          <p>Scan your first book to begin curating your shelves</p>
          <a href="/scanner" class="btn-primary">
            <PlusCircle size="18" />
            <span>Scan a Book</span>
          </a>
        {/if}
      </div>
    {:else}
      {#each filteredBooks as book (book.id || book.Isbn || book.Title)}
        <BookCard {book} />
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

  .tags {
    display: flex;
    overflow-x: auto;
    gap: 0.5rem;
    cursor: pointer;
  }

  .tag-btn {
    background: none;
    border: none;
    padding: 0;
    cursor: pointer;
    display: inline-flex;
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
