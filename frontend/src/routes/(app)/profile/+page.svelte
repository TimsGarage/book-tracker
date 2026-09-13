<script lang="ts">
  import {
    ChevronDown,
    Circle,
    CircleUserRound,
    LogOut,
    ShieldCheck,
  } from "lucide-svelte";
  import { authState } from "$lib/auth.svelte";
  import { handleBack } from "$lib/util";
  import type { NavState } from "../../../lib/nav_helper";
  import { getContext } from "svelte";
  import Input from "../../../components/Input.svelte";
  import { createUser } from "$lib/api";

  function handleLogout() {
    authState.logout();
  }

  function resetNav() {
    navState.heading = "Account";
    navState.showBackButton = false;
    navState.icon = CircleUserRound;
    navState.onBack = handleBack;
    navState.hideBottomNavigation = false;
  }

  let navState = getContext<NavState>("navState");
  if (navState) {
    resetNav();
  }

  let username = $state("");
  let password = $state("");

  let show_admin_options = $state(false);
</script>

<div class="page">
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

  {#if show_admin_options}
    <div class="card">
      <h2>Add User</h2>
      <Input label="Username" bind:value={username}></Input>
      <Input label="Password" bind:value={password}></Input>
      <button onclick={() => createUser(username, password)}>Submit</button>
    </div>
  {:else}
    <button
      style="display: flex; align-items: center; justify-content: space-between; padding-inline: 1.5rem;"
      onclick={() => (show_admin_options = true)}
    >
      Admin Section
      <ChevronDown></ChevronDown>
    </button>
  {/if}
</div>

<style>
  .page {
    padding: 2rem;
    display: flex;
    flex-direction: column;
    gap: 1.5rem;
    height: 100%;
    background-color: var(--main-background);
    overflow-y: auto;
  }

  .card {
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
    color: var(--delete-red);
    border: 2px solid var(--delete-red);
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
</style>
