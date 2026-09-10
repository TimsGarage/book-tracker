<!-- ReadingLoader.svelte -->
<script>
  export let size = "48px";
  export let speed = "3s";
  let className = "";
  export { className as class };
</script>

<div
  class="reading-loader {className}"
  style="--size: {size}; --speed: {speed};"
>
  <svg viewBox="0 0 100 100" xmlns="http://www.w3.org/2000/svg">
    <!-- Abstract text lines -->
    <rect x="0" y="20" width="95" height="8" rx="4" class="line" />
    <rect x="0" y="50" width="95" height="8" rx="4" class="line" />
    <rect x="0" y="80" width="95" height="8" rx="4" class="line" />

    <!-- Reading dot -->
    <g class="dot-wrapper">
      <circle cx="12" cy="24" r="10" class="dot" />
    </g>
  </svg>
</div>

<style>
  .reading-loader {
    width: var(--size);
    height: var(--size);
    display: inline-flex;
    align-items: center;
    justify-content: center;
  }

  svg {
    width: 100%;
    height: 100%;
    overflow: visible;
  }

  .line {
    fill: var(--color);
    opacity: 0.15;
  }

  .dot {
    fill: var(--accent-color);
  }

  .dot-wrapper {
    animation: read var(--speed) infinite cubic-bezier(0.4, 0, 0.2, 1);
  }

  @keyframes read {
    0% {
      transform: translate(-10px, 0);
      opacity: 0;
    }
    10% {
      transform: translate(0, 0);
      opacity: 1;
    }
    25% {
      transform: translate(75px, 0);
      opacity: 1;
    }
    30% {
      transform: translate(85px, 0);
      opacity: 0;
    }

    35% {
      transform: translate(-10px, 30px);
      opacity: 0;
    }
    40% {
      transform: translate(0, 30px);
      opacity: 1;
    }
    55% {
      transform: translate(75px, 30px);
      opacity: 1;
    }
    60% {
      transform: translate(85px, 30px);
      opacity: 0;
    }

    65% {
      transform: translate(-10px, 60px);
      opacity: 0;
    }
    70% {
      transform: translate(0, 60px);
      opacity: 1;
    }
    85% {
      transform: translate(75px, 60px);
      opacity: 1;
    }
    90% {
      transform: translate(85px, 60px);
      opacity: 0;
    }

    100% {
      transform: translate(85px, 60px);
      opacity: 0;
    }
  }
</style>
