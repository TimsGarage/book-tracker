<script lang="ts">
  import { Calendar, ChevronDown, Pencil } from "lucide-svelte";

  interface Props {
    value?: string; // Formatted as "YYYY-MM-DD"
    defaultValue?: string; // Fallback default if value isn't provided
    label?: string; // The floating label text
    placeholder?: string;
    fill?: boolean;
    disabled?: boolean;
    min?: string;
    max?: string;
  }

  let {
    value = $bindable(""),
    defaultValue = "",
    label = "Date",
    placeholder = "Select date",
    fill = false,
    disabled = false,
    min = undefined,
    max = undefined,
  }: Props = $props();

  const id = $props.id();

  // Initialize with defaultValue if no value was passed
  if (!value && defaultValue) {
    value = defaultValue;
  }

  // Formats "YYYY-MM-DD" to a clean localized string (e.g. "Sep 9, 2026")
  let formattedDate = $derived.by(() => {
    if (!value) return "";
    const [year, month, day] = value.split("-").map(Number);
    const date = new Date(year, month - 1, day);
    if (isNaN(date.getTime())) return value;

    return date.toLocaleDateString(undefined, {
      month: "short",
      day: "numeric",
      year: "numeric",
    });
  });
</script>

<div class="datepicker" class:fill class:disabled>
  <!-- Floating label resting on the top border -->
  {#if label}
    <span class="floating-label">{label}</span>
  {/if}

  <!-- Invisible native input spanning the box -->
  <input {id} type="date" bind:value {min} {max} {disabled} />

  <Calendar size="22" strokeWidth="1.5" />

  <label for={id} class="display-text" class:empty={!value}>
    {value ? formattedDate : placeholder}
  </label>

  <Pencil size="18" strokeWidth="1.3" />
</div>

<style>
  .datepicker {
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
    box-sizing: border-box;
    background-color: var(--main-background, transparent);
  }

  .datepicker.fill {
    width: 100%;
  }

  .datepicker.disabled {
    color: var(--text-muted);
    background-color: var(--secondary-background);
    cursor: not-allowed;
  }

  .datepicker.disabled .floating-label {
    visibility: hidden;
  }

  /* Floating label overlapping top border */
  .floating-label {
    position: absolute;
    top: 0;
    left: 0.75rem;
    transform: translateY(-50%);
    background-color: var(--main-background, #fff);
    padding: 0 0.35rem;
    font-size: 0.75rem;
    font-weight: 500;
    color: var(--text-muted, #666);
    line-height: 1;
    pointer-events: none;
    z-index: 2;
  }

  /* Invisible input overlay */
  input[type="date"] {
    opacity: 0;
    position: absolute;
    outline: none;
    border: none;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    cursor: pointer;
    z-index: 1;
  }

  input[type="date"]::-webkit-calendar-picker-indicator,
  input[type="date"]::-webkit-inner-spin-button {
    display: none;
    -webkit-appearance: none;
  }

  /* Visible text inside the box */
  .display-text {
    font-size: 0.95rem;
    font-weight: 400;
    margin-right: auto;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    user-select: none;
  }

  .display-text.empty {
    opacity: 0.6;
  }
</style>
