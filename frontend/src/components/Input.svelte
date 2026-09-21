<script lang="ts">
  let {
    label = "",
    value = $bindable(),
    placeholder = "",
    type = "text",
    icon = null,
    onkeypress = undefined,
  } = $props();

  const id = $props.id();
</script>

<div class="input">
  {#if label != ""}
    <label for={id}>{label}</label>
  {/if}
  {#if icon}
    <span class="icon">
      {@render icon()}
    </span>
  {/if}
  <input {onkeypress} {id} bind:value {type} {placeholder} />
</div>

<style>
  .input {
    position: relative;
    display: flex;
    flex-direction: column;
  }

  input {
    height: 3.5rem;
    outline: none;
    border: none;
    background-color: var(--secondary-background);
    border-radius: 0.5rem;
    padding: 0.75rem 1rem;
  }

  label {
    position: absolute;
    top: 32.5%;
    left: 20px;
    scale: 1;
    transition: all 0.125s ease-in-out;
    opacity: 0.8;
  }

  .input:has(> input:focus) label,
  .input:has(> input:not(:placeholder-shown)) label {
    top: -18%;
    left: 10px;
    scale: 0.9;
    opacity: 1;
  }

  .input:has(.icon) input {
    padding-left: 3.5rem;
  }

  .icon {
    position: absolute;
    width: 3.5rem;
    height: 3.5rem;
    left: 0;
    top: 0;
    display: flex;
    justify-content: center;
    align-items: center;
  }
</style>
