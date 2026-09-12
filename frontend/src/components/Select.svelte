<script lang="ts">
  import { ChevronDown } from "lucide-svelte";

  let { options, value = $bindable(), fill = false } = $props();

  const id = $props.id();

  // svelte-ignore state_referenced_locally
  let _value = $state(value ?? options[0].value ?? "");
  let selected = $derived(options.filter((e: any) => e.value == _value)[0]);

  $effect(() => {
    value = _value;
  });
</script>

<div class="select" class:fill>
  <select {id} bind:value={_value}>
    {#each options as option}
      <option value={option.value}>
        {option.title}
      </option>
    {/each}
  </select>

  <selected.icon size="24" strokeWidth="1.5"></selected.icon>
  <label for={id}>{selected.title}</label>
  <ChevronDown size="24" strokeWidth="1.3" />
</div>

<style>
  select {
    opacity: 0;
    position: absolute;
    outline: none;
    border: none;
    top: 2.5%;
    left: 2.5%;
    width: 95%;
    height: 95%;
  }

  .select {
    position: relative;
    height: 3.5rem;
    min-width: 8rem;
    max-width: 100%;
    border-radius: 0.5rem;
    border: 1px solid var(--outline);
    padding: 1rem 0.75rem;
    display: flex;
    align-items: center;
    gap: 0.75rem;
    background-color: var(--main-background, transparent);
  }

  .select.fill {
    width: 100%;
  }

  label {
    font-size: 0.95rem;
    font-weight: 400;
    margin-right: auto;
  }
</style>
