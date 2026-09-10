<script lang="ts">
  import { Star } from "lucide-svelte";
  import Chip from "./Chip.svelte";
  import Loader from "./Loader.svelte";
  let { book } = $props();

  let loading_cover = $state(true);
</script>

<div class="book-preview">
  {#if loading_cover}
    <div class="thumbnail">
      <Loader size="32px" />
    </div>
  {/if}
  <img
    src={book.ThumbnailLink}
    onload={() => (loading_cover = false)}
    alt="cover"
    class="thumbnail"
  />

  <div class="publisher">
    <p>{book.Publisher}</p>
    <p>Unread</p>
  </div>

  <div class="main-info">
    <h2>{book.Title}</h2>
    <h6>{book.Author}</h6>
  </div>

  <div class="additional-info">
    <p>
      Pages: {book.Pages} - {book.Release} -
      <span>3/5<Star size="16" color="var(--accent-color)" /></span>
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
