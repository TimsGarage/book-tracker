<script lang="ts">
  import Loader from "./Loader.svelte";
  import type { Book } from "$lib/types";
  let { book, onclick }: { book: Book; onclick: any } = $props();

  let loading_cover = $state(true);
</script>

<!-- svelte-ignore a11y_click_events_have_key_events -->
<!-- svelte-ignore a11y_no_static_element_interactions -->
<div class="book-preview" {onclick}>
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
      <span>No Cover</span>
    </div>
  {/if}

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
</div>

<style>
  .book-preview {
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

  .book-preview > .thumbnail {
    grid-area: thumbnail;
    height: 100%;
    width: 100%;
    object-fit: contain;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .book-preview .publisher {
    grid-area: publisher;
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .book-preview .main-info {
    grid-area: main-info;
    align-self: flex-start;
  }

  .book-preview .additional-info {
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
  }
</style>
