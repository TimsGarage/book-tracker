<script lang="ts">
  import { goto } from "$app/navigation";
  import BookCard from "./BookCard.svelte";

  let { books } = $props();
</script>

<div class="grid-view">
  {#each books as book (book.id)}
    <BookCard onclick={() => goto("/" + book.id)} {book} reduced />
  {/each}
</div>

<style>
  .grid-view {
    width: 100%;
    height: 100%;
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    justify-items: center;
    gap: 1rem;
    padding: 1rem;
  }
</style>
