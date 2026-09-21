<script lang="ts">
  import { goto } from "$app/navigation";
  import BookCard from "./BookCard.svelte";

  let { books } = $props();
</script>

{#each books as book (book.id)}
  <BookCard onclick={() => goto("/" + book.id)} {book} />
{/each}
