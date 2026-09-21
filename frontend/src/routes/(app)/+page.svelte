<script lang="ts">
  import { getContext, onMount } from "svelte";
  import BookCard from "../../components/BookCard.svelte";
  import { type NavState } from "../../lib/nav_helper";
  import Input from "../../components/Input.svelte";
  import Loader from "../../components/Loader.svelte";
  import { Search, BookOpen, CirclePlus, Icon, Library } from "lucide-svelte";
  import Chip from "../../components/Chip.svelte";
  import { fetchMyBooks } from "$lib/api";
  import type {
    Book,
    BookOwnershipStatus,
    BookReadingStatus,
  } from "$lib/types";
  import BookPage from "../../components/BookPreview.svelte";
  import { handleBack, owning_options } from "$lib/util";
  import Select from "../../components/Select.svelte";
  import { goto } from "$app/navigation";
  import ListView from "../../components/ListView.svelte";
  import GridView from "../../components/GridView.svelte";

  let navState = getContext<NavState>("navState");
  if (navState) {
    navState.heading = "Library";
    navState.showBackButton = false;
    navState.onBack = handleBack;
    navState.icon = Library;
    navState.hideBottomNavigation = false;
  }

  let books = $state<Book[]>([]);
  let isLoading = $state(true);
  let errorMessage = $state("");
  let gridView = $state(true);
  let searchQuery = $state("");
  let activeOwnershipFilter = $state<"all" | BookOwnershipStatus>("all");
  let activeReadingStatusFilter = $state<"all" | BookReadingStatus>("all");

  async function loadBooks() {
    isLoading = true;
    errorMessage = "";
    try {
      books = await fetchMyBooks();
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

  // 1. Filter books ONLY by search query
  const searchMatchedBooks = $derived(
    books.filter((b: Book) => {
      let searchMatched = false;
      const q = searchQuery.toLowerCase().trim();
      if (!q) searchMatched = true;
      searchMatched =
        (b.title || "").toLowerCase().includes(q) ||
        (b.author || "").toLowerCase().includes(q) ||
        (b.isbn || "").toLowerCase().includes(q);

      let ownershipFilterMatched = false;
      if (activeOwnershipFilter === "all") {
        ownershipFilterMatched = b.ownership_status !== "wishlist";
      } else {
        ownershipFilterMatched = b.ownership_status === activeOwnershipFilter;
      }

      return ownershipFilterMatched && searchMatched;
    }),
  );

  // 2. Derive counts directly from searchMatchedBooks
  const unreadCount = $derived(
    searchMatchedBooks.filter((b) => b.reading_status === "unread").length,
  );
  const readCount = $derived(
    searchMatchedBooks.filter((b) => b.reading_status === "read").length,
  );
  const readingCount = $derived(
    searchMatchedBooks.filter((b) => b.reading_status === "reading").length,
  );

  // 3. Derive final list by applying the active status tab filter
  const filteredBooks = $derived(
    searchMatchedBooks.filter((b) => {
      let readingFilterMatched = false;
      if (activeReadingStatusFilter === "all") {
        readingFilterMatched = true;
      } else {
        readingFilterMatched = b.reading_status === activeReadingStatusFilter;
      }

      return readingFilterMatched;
    }),
  );
</script>

<div class="page">
  <div class="header">
    <div
      class="group"
      style="display: grid; grid-template-columns: 60% 40%; gap: .25rem;"
    >
      <Input placeholder="Search for a Book" bind:value={searchQuery}>
        {#snippet icon()}
          <Search size="24" color="var(--text)" />
        {/snippet}
      </Input>
      <Select
        bind:value={activeOwnershipFilter}
        options={[
          {
            title: "All Books",
            value: "all",
          },
          ...owning_options.filter((e) => !(e.value == "wishlist")),
        ]}
      />
    </div>

    <div class="tags">
      <Chip
        onclick={() => (activeReadingStatusFilter = "all")}
        text="All Books ({filteredBooks.length})"
        selected={activeReadingStatusFilter === "all"}
      />
      <Chip
        onclick={() => (activeReadingStatusFilter = "unread")}
        text="Unread ({unreadCount})"
        selected={activeReadingStatusFilter === "unread"}
      />
      <Chip
        text="Read ({readCount})"
        selected={activeReadingStatusFilter === "read"}
        onclick={() => (activeReadingStatusFilter = "read")}
      />
      <Chip
        text="Reading ({readingCount})"
        selected={activeReadingStatusFilter === "reading"}
        onclick={() => (activeReadingStatusFilter = "reading")}
      />
    </div>

    <div class="view">
      <button class:active={!gridView} onclick={() => (gridView = false)}
        >List</button
      >
      <button class:active={gridView} onclick={() => (gridView = true)}
        >Grid</button
      >
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
            <CirclePlus size="18" />
            <span>Scan a Book</span>
          </a>
        {/if}
      </div>
    {:else if gridView}
      <GridView books={filteredBooks} />
    {:else}
      <ListView books={filteredBooks} />
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
    padding-bottom: 0px;
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

  .view {
    width: 100%;
    display: flex;
  }

  .view button {
    height: 2rem;
    min-height: unset;
    border-collapse: collapse;
    border-radius: 1rem;
  }

  .view button.active {
    background-color: var(--accent-color);
    color: var(--text-reversed);
    border-color: var(--accent-color);
  }

  .view button:nth-of-type(1) {
    border-top-right-radius: 0px;
    border-bottom-right-radius: 0px;
  }

  .view button:nth-of-type(2) {
    border-top-left-radius: 0px;
    border-bottom-left-radius: 0px;
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
    height: 90%;
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
