<script lang="ts">
  import { goto } from "$app/navigation";
  import Checkbox from "../../../components/Checkbox.svelte";
  import Input from "../../../components/Input.svelte";
</script>

<div class="page">
  <div class="content">
    <div class="header">
      <p>Member Sign In</p>
      <h1>Welcome back to your library</h1>
      <p>Sign in to catalog your physical books and curate your home shelves</p>
    </div>

    <div class="form">
      <Input label="Username" />
      <Input label="Password" type="password" />

      <Checkbox label={"Remember me"} />

      <button onclick={() => goto("/")} class="Primary">Sign In</button>
    </div>
  </div>
</div>

<style>
  .page {
    height: 100%;
    width: 100%;
    background-color: var(--main-background);
    display: flex;
    justify-content: center;
    align-items: center;
  }

  .content {
    width: 80%;
  }

  .header {
    text-align: center;
    margin-bottom: 1.5rem;
  }

  .form {
    padding: 1rem;
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }
</style>
