<script lang="ts">
  import { goto } from "$app/navigation";
  import Checkbox from "../../../components/Checkbox.svelte";
  import Input from "../../../components/Input.svelte";
  import Loader from "../../../components/Loader.svelte";
  import { authState } from "$lib/auth.svelte";

  let username = $state("");
  let password = $state("");
  let rememberMe = $state(true);
  let errorMessage = $state("");
  let isSubmitting = $state(false);

  async function handleLogin(e?: Event) {
    if (e) e.preventDefault();
    errorMessage = "";

    if (!username.trim()) {
      errorMessage = "Please enter your username";
      return;
    }
    if (!password) {
      errorMessage = "Please enter your password";
      return;
    }

    isSubmitting = true;
    try {
      await authState.login(
        {
          username: username.trim(),
          password,
        },
        rememberMe,
      );
      goto("/");
    } catch (err: any) {
      errorMessage = err.message || "Invalid username or password";
    } finally {
      isSubmitting = false;
    }
  }
</script>

<div class="page">
  <div class="content">
    <div class="header">
      <p>Member Sign In</p>
      <h1>Welcome back to your library</h1>
      <p>Sign in to catalog your physical books and curate your home shelves</p>
    </div>

    <form class="form" onsubmit={handleLogin}>
      {#if errorMessage}
        <div class="error-banner" role="alert">
          {errorMessage}
        </div>
      {/if}

      <Input label="Username" bind:value={username} />
      <Input label="Password" type="password" bind:value={password} />

      <Checkbox label={"Remember me"} bind:value={rememberMe} />

      <button type="submit" class="Primary" disabled={isSubmitting}>
        {#if isSubmitting}
          <span class="btn-loader">
            <Loader size="20px" />
            <span>Signing in...</span>
          </span>
        {:else}
          Sign In
        {/if}
      </button>
    </form>
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
    max-width: 420px;
  }

  .header {
    text-align: center;
    margin-bottom: 1.5rem;
  }

  .header h1 {
    font-size: 1.75rem;
    margin-block: 0.5rem;
  }

  .header p {
    color: var(--text-muted);
    font-size: 0.9rem;
  }

  .form {
    padding: 1rem;
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }

  .error-banner {
    background-color: #fee2e2;
    color: #b91c1c;
    border: 1px solid #f87171;
    border-radius: 0.5rem;
    padding: 0.75rem 1rem;
    font-size: 0.875rem;
    font-weight: 500;
  }

  button.Primary {
    height: 3.25rem;
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    font-weight: 600;
  }

  button.Primary:disabled {
    opacity: 0.7;
    cursor: not-allowed;
  }

  .btn-loader {
    display: inline-flex;
    align-items: center;
    gap: 0.5rem;
  }
</style>
