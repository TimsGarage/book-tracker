<script lang="ts">

    let {
        label,
        value = $bindable(),
    } = $props();

    const id = $props.id();

</script>

<div class="input">
    <label for={id}>{label}</label>
    <input {id} placeholder="" bind:checked={value} type="checkbox"/>
</div>

<style>

.input {
    display: flex;
    align-items: center;
    justify-content: space-between;
}

input {
    height: 1rem;
    outline: none;
    border: none;
    background-color: var(--secondary-background);
    border-radius: .25rem;
    padding: .5rem 1rem;
}

label {
    scale: .9;
    opacity: .8;

}

</style>