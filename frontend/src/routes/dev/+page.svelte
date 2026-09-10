<script lang="ts">
  import {
    CircleUserRound,
    LogOut,
    ArrowLeft,
    ShieldCheck,
  } from "lucide-svelte";
  import Loader from "../../components/Loader.svelte";
  import { authState } from "$lib/auth.svelte";
  import { goto } from "$app/navigation";

  function handleLogout() {
    authState.logout();
  }
</script>

<div class="page">
  <div class="header">
    <button class="icon-btn" onclick={() => goto("/")}>
      <ArrowLeft size="24" />
    </button>
    <h1>Account & Settings</h1>
  </div>

  <div class="card account-card">
    <div class="avatar">
      <CircleUserRound size="64" color="var(--accent-color)" />
    </div>

    {#if authState.user}
      <div class="user-details">
        <h2>{authState.user.username}</h2>
        <p class="user-id">User ID: #{authState.user.id}</p>
        <span class="badge">
          <ShieldCheck size="14" />
          <span>Authenticated</span>
        </span>
      </div>
    {:else}
      <div class="user-details">
        <p>Not signed in</p>
      </div>
    {/if}

    <button class="logout-btn" onclick={handleLogout}>
      <LogOut size="18" />
      <span>Sign Out</span>
    </button>
  </div>

  <div class="card components-card">
    <h3>Design System Preview</h3>
    <div class="loader-demo">
      <Loader size="48px" />
    </div>
  </div>
</div>

<style>
  .page {
    padding: 2rem 1.5rem;
    display: flex;
    flex-direction: column;
    gap: 1.5rem;
    height: 100%;
    background-color: var(--main-background);
    overflow-y: auto;
  }

  .header {
    display: flex;
    align-items: center;
    gap: 1rem;
    padding-top: 1rem;
  }

  .header h1 {
    font-size: 1.5rem;
    margin: 0;
  }

  .icon-btn {
    background: none;
    border: none;
    color: var(--text);
    cursor: pointer;
    padding: 0.5rem;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .card {
    background-color: var(--secondary-background);
    border: 1px solid var(--outline);
    border-radius: 1rem;
    padding: 1.5rem;
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }

  .account-card {
    align-items: center;
    text-align: center;
  }

  .user-details h2 {
    margin: 0 0 0.25rem 0;
    font-size: 1.5rem;
  }

  .user-id {
    color: var(--text-muted);
    font-size: 0.875rem;
    margin: 0 0 0.5rem 0;
  }

  .badge {
    display: inline-flex;
    align-items: center;
    gap: 0.35rem;
    background-color: rgba(34, 197, 94, 0.15);
    color: #16a34a;
    font-size: 0.75rem;
    font-weight: 600;
    padding: 0.25rem 0.6rem;
    border-radius: 1rem;
  }

  .logout-btn {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    gap: 0.5rem;
    background-color: #fee2e2;
    color: #dc2626;
    border: 1px solid #fca5a5;
    padding: 0.75rem 1.5rem;
    border-radius: 0.5rem;
    cursor: pointer;
    font-weight: 600;
    width: 100%;
    margin-top: 0.5rem;
    transition: background-color 0.15s ease;
  }

  .logout-btn:hover {
    background-color: #fecaca;
  }

  .components-card h3 {
    margin: 0;
    font-size: 1rem;
    color: var(--text-muted);
  }

  .loader-demo {
    display: flex;
    justify-content: center;
    padding: 1rem;
  }
</style>
