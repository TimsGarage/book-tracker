import { browser } from "$app/environment";

export function createPersistentState<T>(key: string, initialValue: T) {
  let stored: T = initialValue;

  if (browser) {
    const raw = localStorage.getItem(key);
    if (raw !== null) {
      try {
        stored = JSON.parse(raw);
      } catch {
        stored = initialValue;
      }
    }
  }

  let value = $state<T>(stored);

  $effect(() => {
    if (browser) {
      localStorage.setItem(key, JSON.stringify(value));
    }
  });

  return {
    get current() {
      return value;
    },
    set current(v: T) {
      value = v;
    },
  };
}