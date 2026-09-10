import { redirect } from '@sveltejs/kit';
import { authState } from '$lib/auth.svelte';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ fetch }) => {
  const user = await authState.verifyToken(fetch);
  if (!user) {
    redirect(307, '/login');
  }
  return { user };
};

