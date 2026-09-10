import { browser } from '$app/environment';
import { goto } from '$app/navigation';
import type { User, AuthResponse, LoginCredentials } from './types';

export const API_BASE = import.meta.env.VITE_API_URL || '';

class AuthStore {
  user = $state<User | null>(null);
  token = $state<string | null>(null);
  isLoading = $state<boolean>(false);
  isAuthenticated = $derived(!!this.token && !!this.user);

  constructor() {
    if (browser) {
      this.initFromStorage();
    }
  }

  initFromStorage() {
    try {
      const token = localStorage.getItem('token') || sessionStorage.getItem('token');
      const userStr = localStorage.getItem('user') || sessionStorage.getItem('user');

      if (token) {
        this.token = token;
      }
      if (userStr) {
        this.user = JSON.parse(userStr);
      }
    } catch (e) {
      console.warn('Failed to parse user from storage:', e);
      this.clearStorage();
    }
  }

  setSession(token: string, user: User, remember: boolean = true) {
    this.token = token;
    this.user = user;

    if (browser) {
      this.clearStorage();
      const storage = remember ? localStorage : sessionStorage;
      storage.setItem('token', token);
      storage.setItem('user', JSON.stringify(user));
    }
  }

  clearStorage() {
    if (browser) {
      localStorage.removeItem('token');
      localStorage.removeItem('user');
      sessionStorage.removeItem('token');
      sessionStorage.removeItem('user');
    }
    this.token = null;
    this.user = null;
  }

  getToken(): string | null {
    if (this.token) return this.token;
    if (browser) {
      return localStorage.getItem('token') || sessionStorage.getItem('token');
    }
    return null;
  }

  logout() {
    this.clearStorage();
    if (browser) {
      goto('/login');
    }
  }

  async login(credentials: LoginCredentials, remember: boolean = true): Promise<User> {
    this.isLoading = true;
    try {
      const jsonBody = JSON.stringify(credentials);
      let b64Payload = "";
      try {
        b64Payload = btoa(unescape(encodeURIComponent(jsonBody)));
      } catch (e) {
        console.warn("Could not base64 encode payload:", e);
      }

      // Append query parameters as fallback for Android WebView body stripping
      const baseUrl = `${API_BASE}/api/v1/auth/login`;
      const url = new URL(
        baseUrl,
        typeof window !== "undefined" ? window.location.origin : "http://localhost"
      );
      url.searchParams.set("username", credentials.username);
      url.searchParams.set("password", credentials.password);

      const headers: Record<string, string> = {
        "Content-Type": "application/json",
      };
      if (b64Payload) {
        headers["X-Payload"] = b64Payload;
      }

      const res = await fetch(url.toString(), {
        method: "POST",
        headers,
        body: jsonBody,
      });

      const data = await res.json();

      if (!res.ok) {
        throw new Error(data.error || "Invalid credentials");
      }

      const authResp = data as AuthResponse;
      this.setSession(authResp.token, authResp.user, remember);
      return authResp.user;
    } finally {
      this.isLoading = false;
    }
  }

  async verifyToken(customFetch?: typeof fetch): Promise<User | null> {
    const fetchFn = customFetch || fetch;
    const currentToken = this.getToken();

    if (!currentToken) {
      this.clearStorage();
      return null;
    }

    try {
      const res = await fetchFn(`${API_BASE}/api/v1/auth/me`, {
        headers: {
          Authorization: `Bearer ${currentToken}`,
        },
      });

      if (!res.ok) {
        // Token is invalid, expired, or rejected
        this.clearStorage();
        return null;
      }

      const user = (await res.json()) as User;
      this.token = currentToken;
      this.user = user;

      if (browser) {
        const isLocalStorage = !!localStorage.getItem("token");
        const storage = isLocalStorage ? localStorage : sessionStorage;
        storage.setItem("token", currentToken);
        storage.setItem("user", JSON.stringify(user));
      }

      return user;
    } catch (err) {
      console.error("Failed to verify token:", err);
      return null;
    }
  }

  async authFetch(
    input: string | URL | Request,
    init?: RequestInit,
    customFetch?: typeof fetch
  ): Promise<Response> {
    const fetchFn = customFetch || fetch;
    const token = this.getToken();

    const headers = new Headers(init?.headers || {});
    if (token) {
      headers.set("Authorization", `Bearer ${token}`);
    }

    if (init?.body && typeof init.body === "string") {
      if (!headers.has("Content-Type")) {
        headers.set("Content-Type", "application/json");
      }
      try {
        const b64 = btoa(unescape(encodeURIComponent(init.body)));
        headers.set("X-Payload", b64);
      } catch (e) {
        console.warn("Failed to encode X-Payload:", e);
      }
    }

    const response = await fetchFn(input, {
      ...init,
      headers,
    });

    if (response.status === 401) {
      this.logout();
    }

    return response;
  }
}

export const authState = new AuthStore();

