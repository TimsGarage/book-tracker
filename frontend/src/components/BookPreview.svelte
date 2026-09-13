<script lang="ts">
  import Select from "./Select.svelte";
  import Datepicker from "./Datepicker.svelte";
  import { getTodayString, owning_options, read_options } from "$lib/util";
  import type { Book, LookupBook } from "$lib/types";
  import { Save, Trash } from "lucide-svelte";
  import Loader from "./Loader.svelte";

  let {
    book = $bindable(),
    editMode = false,
    deleteCallback = ({}) => {},
    saveCallback = ({}) => {},
  }: {
    book: Book;
    editMode?: boolean;
    deleteCallback?: (book: Book) => void;
    saveCallback?: (book: Book) => void;
  } = $props();

  function onDelete() {
    deleteCallback(myBook);
  }

  function onSave() {
    saveCallback(myBook);
  }

  let myBook = $state(book);
  myBook.owned_since = myBook.owned_since ?? getTodayString();

  let loading_cover = $state(true);
</script>

<div class="bookpreview" style="margin-bottom: 1.5rem;">
  <span class="decoration-1"></span>
  <span class="decoration-2"></span>

  {#if editMode}
    <button class="deleteButton" onclick={onDelete}>
      <Trash></Trash>
    </button>
    <button class="saveButton" onclick={onSave}>
      <Save></Save>
    </button>
  {/if}

  <div class="thumbnail">
    {#if loading_cover}
      <Loader size="64px"></Loader>
    {/if}
    {#if myBook.thumbnail_link}
      <img
        onload={() => {
          loading_cover = false;
        }}
        onerror={() => (loading_cover = false)}
        src={myBook.thumbnail_link}
        alt={myBook.title || "cover"}
      />
    {/if}
  </div>

  <h1 class="title">{myBook.title}</h1>
  <h3 class="author">
    {myBook.author}
  </h3>
  <p class="additional-info">
    {myBook.pages ? `${myBook.pages} Pages` : ""}
  </p>

  <p
    style="width: 100%; font-size: .8rem; color: var(--text-muted); font-weight: 600; border-bottom: 1px solid var(--outline); margin-top: .75rem; padding-block: .25rem;"
  >
    Catalogue options
  </p>

  <div class="button-group">
    <Select
      options={owning_options}
      bind:value={myBook.ownership_status}
      fill
    />
    <Datepicker
      label="Purchased"
      bind:value={myBook.owned_since}
      fill
      disabled={myBook.ownership_status != "owned"}
    />
  </div>
  <Select options={read_options} bind:value={myBook.reading_status} fill />
</div>

<style>
  .bookpreview {
    z-index: 1;
    position: relative;
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 0.5rem;
    padding: 1.5rem;
    height: 100%;
    width: 100%;
    overflow: hidden;
  }

  .bookpreview > * {
    z-index: 1;
  }

  .decoration-1,
  .decoration-2 {
    z-index: 0;
    content: "";
    position: absolute;
    width: 500px;
    height: 500px;
    border-radius: 40px;
    background-color: var(--accent-color);
    rotate: 35deg;
  }

  .decoration-1 {
    left: -470px;
    top: -80px;
  }

  .decoration-2 {
    right: -480px;
    bottom: -90px;
  }

  .deleteButton,
  .saveButton {
    position: absolute;
    z-index: 5;
    right: 1.5rem;
    top: 1.5rem;
    min-width: unset;
    min-height: unset;
    height: 64px;
    width: 64px;
    border-radius: 100%;
    background-color: var(--accent-color);
    display: flex;
    align-items: center;
    justify-content: center;
    color: var(--main-background);
  }

  .deleteButton {
    color: var(--delete-red);
    background-color: transparent;
    border: 2px solid var(--delete-red);
  }

  .saveButton {
    /* Claculate deleteButtonTop + deletebuttonsize + padding */
    top: calc(1.5rem + 64px + 0.5rem);
  }

  .thumbnail {
    max-width: 250px;
    min-height: 300px;
    flex-shrink: 1;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
  }

  .thumbnail img {
    object-fit: contain;
    max-width: 100%;
  }

  .title {
    text-align: center;
    margin-inline: 1rem;
    font-size: 2.5rem;
  }

  .author,
  .additional-info {
    text-align: center;
    margin-inline: 1rem;
    opacity: 0.75;
  }

  .button-group {
    display: flex;
    justify-content: space-evenly;
    gap: 0.5rem;
    width: 100%;
  }
</style>
