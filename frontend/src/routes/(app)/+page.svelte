<script lang="ts">
  import { getContext } from "svelte";
  import BookCard from "../../components/BookCard.svelte";
  import type { TopNavState } from "../../components/types";
  import Input from "../../components/Input.svelte";
  import { Search } from "lucide-svelte";
  import Chip from "../../components/Chip.svelte";

  let topNavState = getContext<TopNavState>("topNavState");
  topNavState.heading = "Home";
  topNavState.showBackButton = false;

  const test_data = [
    {
      Isbn: "9783792000274",
      Title: "Der Kleine Prinz",
      Author: "Antoine de Saint-Exupéry",
      Publisher: "Verlag",
      Release: 1979,
      Description: "",
      ThumbnailLink: "https://covers.openlibrary.org/b/id/1027413-L.jpg",
      Pages: 94,
    },
    {
      Isbn: "9781546154419",
      Title: "Harry Potter and the Goblet of Fire",
      Author: "J. K. Rowling",
      Publisher: "Verlag",
      Release: 1979,
      Description: "",
      ThumbnailLink: "https://covers.openlibrary.org/b/id/15227222-L.jpg",
      Pages: 608,
    },
    {
      Isbn: "9783743200111",
      Title: "Elanus",
      Author: "Ursula Poznanski",
      Publisher: "Verlag",
      Release: 1979,
      Description: "",
      ThumbnailLink: "https://covers.openlibrary.org/b/id/14570203-L.jpg",
      Pages: 412,
    },
  ];
</script>

<div class="page">
  <div class="header">
    <Input placeholder="Search for a Book">
      {#snippet icon()}
        <Search size="24" color="var(--text)" />
      {/snippet}
    </Input>

    <div class="tags">
      <Chip text="All Books (25)" selected />
      <Chip text="Ungelesen (12)" />
      <Chip text="Gelesen (6)" />
      <Chip text="Am Lesen (11)" />
    </div>
  </div>

  {#each test_data as book}
    <BookCard {book} />
  {/each}
</div>

<style>
  .page {
    height: 100%;
    width: 100%;
    position: relative;
    display: flex;
    flex-direction: column;
    background-color: var(--main-background);
  }

  .header {
    display: flex;
    flex-direction: column;
    gap: 1rem;
    padding: 2rem 1rem;
    padding-bottom: 1rem;
    border-bottom: 1px solid var(--outline);
  }

  .tags {
    display: flex;
    overflow-x: auto;
    gap: 0.5rem;
  }
</style>
