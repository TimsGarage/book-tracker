<script lang="ts">

    let {
        selected = false,
        text = ""
    } = $props();

</script>

<span class:selected>
    <p>{text}</p>
</span>

<style>

span {
    display: flex;
    justify-content: center;
    align-items: center;
    padding: .75rem 1.25rem;
    border-radius: 2rem;
    min-width: max-content;
    background-color: var(--secondary-background);
    border: 1px solid var(--outline);
}

span p {
    text-wrap: nowrap;
}

span.selected {
    background-color: var(--accent-color);
    color: var(--text-reversed);
    border-color: var(--accent-outline);
}
    
</style>