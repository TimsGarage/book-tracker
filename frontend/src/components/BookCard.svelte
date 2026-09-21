<script lang="ts">
  import Loader from "./Loader.svelte";
  import type { Book } from "$lib/types";
  let {
    book,
    onclick,
    reduced = false,
  }: { book: Book; onclick: any; reduced?: boolean } = $props();

  let loading_cover = $state(true);
</script>

<!-- svelte-ignore a11y_click_events_have_key_events -->
<!-- svelte-ignore a11y_no_static_element_interactions -->
<div class="book-card" class:reduced {onclick}>
  {#if loading_cover && book.thumbnail_link}
    <div class="thumbnail">
      <Loader size="32px" />
    </div>
  {/if}
  {#if book.thumbnail_link}
    <img
      src={book.thumbnail_link}
      onload={() => (loading_cover = false)}
      onerror={() => (loading_cover = false)}
      alt={book.title || "cover"}
      class="thumbnail"
    />
  {:else}
    <div class="thumbnail no-cover">
      <span>{book.title}</span>
    </div>
  {/if}

  {#if !reduced}
    <div class="publisher">
      <p>{book.publisher}</p>
      <p>{book.ownership_status}</p>
    </div>

    <div class="main-info">
      <h2>{book.title}</h2>
      <h6>{book.author}</h6>
    </div>

    <div class="additional-info">
      <p>
        Pages: {book.pages} - {book.release} - {book.reading_status}
        <!-- <span>3/5<Star size="16" color="var(--accent-color)" /></span> -->
      </p>
    </div>
  {/if}
</div>

<style>
  .book-card {
    --height: 150px;
    width: 100%;
    height: var(--height);
    border-bottom: 1px solid var(--outline);
    padding: 1rem;
    display: grid;
    grid-template-areas:
      "thumbnail publisher"
      "thumbnail main-info"
      "thumbnail additional-info";
    grid-template-columns: 100px 1fr;
    grid-template-rows: min-content 1fr min-content;
    column-gap: 0.5rem;
  }

  .book-card.reduced {
    display: flex;
    background-color: var(--secondary-background);
    border: none;
    border-radius: 0.25rem;
    padding: 0.5rem;
  }

  .book-card > .thumbnail {
    grid-area: thumbnail;
    height: 100%;
    width: 100%;
    object-fit: contain;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .book-card .publisher {
    grid-area: publisher;
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .book-card .main-info {
    grid-area: main-info;
    align-self: flex-start;
  }

  .book-card .additional-info {
    grid-area: additional-info;
    height: fit-content;
    align-self: flex-end;
    display: flex;
    align-items: center;
  }

  p,
  .publisher {
    font-size: 0.8rem;
    color: var(--text-muted);
  }

  h2 {
    line-height: 1.6rem;
    display: -webkit-box;
    -webkit-box-orient: vertical;
    -webkit-line-clamp: 2;
    line-clamp: 2; /* Standard property for future compatibility */
    overflow: hidden;
  }
</style>
