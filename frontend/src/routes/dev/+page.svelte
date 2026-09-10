<script lang="ts">
  import {
    Book,
    BookCheck,
    Bookmark,
    BookMarked,
    Info,
    Lock,
    ScrollText,
  } from "lucide-svelte";
  import Select from "../../components/Select.svelte";
  import Input from "../../components/Input.svelte";
  import Datepicker from "../../components/Datepicker.svelte";
  import Loader from "../../components/Loader.svelte";

  const options = [
    {
      title: "To be Read",
      value: "tbr",
      icon: BookMarked,
    },
    {
      title: "Wishlist",
      value: "wishlist",
      icon: ScrollText,
    },
    {
      title: "Read",
      value: "read",
      icon: BookCheck,
    },
  ];
</script>

<div class="page">
  <h1>test</h1>

  <Loader size="48px" />
</div>

<style>
  .page {
    padding: 3rem;
    padding-top: 3rem;
  }
</style>
