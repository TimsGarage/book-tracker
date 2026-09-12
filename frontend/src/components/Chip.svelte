<script lang="ts">
  let {
    selected = false,
    text = "",
    onclick = undefined,
  }: {
    selected?: boolean;
    text?: string;
    onclick?: (event: MouseEvent) => void;
  } = $props();
</script>

<button {onclick} class:selected>
  <p>{text}</p>
</button>

<style>
  button {
    min-height: unset;
    height: fit-content;
    display: flex;
    justify-content: center;
    align-items: center;
    padding: 0.75rem 1.25rem;
    border-radius: 2rem;
    min-width: max-content;
    background-color: var(--secondary-background);
    border: 1px solid var(--outline);
  }

  button p {
    text-wrap: nowrap;
  }

  button.selected {
    background-color: var(--accent-color);
    color: var(--text-reversed);
    border-color: var(--accent-outline);
  }
</style>
