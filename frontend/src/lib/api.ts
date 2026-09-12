import { authState, API_BASE } from './auth.svelte';
import type { Book, LookupBook } from './types';

export async function createUser(username: string, password: string, customFetch?: typeof fetch): Promise<any> {
  const baseUrl = `${API_BASE}/api/v1/auth/register`;
  const url = new URL(baseUrl, typeof window !== "undefined" ? window.location.origin : "http://localhost");
  url.searchParams.set("username", username);
  url.searchParams.set("password", password);

  const res = await authState.authFetch(
    url.toString(),
    {
      method: "POST",
      body: JSON.stringify({
        username,
        password,
      }),
    },
    customFetch
  );
  if (!res.ok) {
    const errorData = await res.json().catch(() => ({}));
    throw new Error(errorData.error || `Failed to create user: ${res.statusText}`);
  }
  const data = await res.json();
  return data.data;
}

export async function fetchOwnedBooks(customFetch?: typeof fetch): Promise<Book[]> {
  const res = await authState.authFetch(`${API_BASE}/api/v1/books/owned`, {}, customFetch);
  if (!res.ok) {
    throw new Error(`Failed to fetch books: ${res.statusText}`);
  }
  const data = await res.json();
  return data.data;
}

export async function fetchReadBooks(customFetch?: typeof fetch): Promise<Book[]> {
  const res = await authState.authFetch(`${API_BASE}/api/v1/books/read`, {}, customFetch);
  if (!res.ok) {
    throw new Error(`Failed to fetch books: ${res.statusText}`);
  }
  const data = await res.json();
  return data.data;
}

export async function fetchWishlistBooks(customFetch?: typeof fetch): Promise<Book[]> {
  const res = await authState.authFetch(`${API_BASE}/api/v1/books/wishlist`, {}, customFetch);
  if (!res.ok) {
    throw new Error(`Failed to fetch books: ${res.statusText}`);
  }
  const data = await res.json();
  return data.data;
}

export async function fetchBookById(id: number, customFetch?: typeof fetch): Promise<Book> {
  const res = await authState.authFetch(`${API_BASE}/api/v1/books/${id}`, {}, customFetch);
  if (!res.ok) {
    throw new Error(`Failed to fetch book: ${res.statusText}`);
  }
  const data = await res.json();
  return data.data;
}

export async function createBook(book: Book, customFetch?: typeof fetch): Promise<Book> {
  const baseUrl = `${API_BASE}/api/v1/books`;
  const url = new URL(baseUrl, typeof window !== "undefined" ? window.location.origin : "http://localhost");
  if (book.isbn) url.searchParams.set("isbn", book.isbn);
  if (book.title) url.searchParams.set("title", book.title);
  if (book.author) url.searchParams.set("author", book.author);
  if (book.description) url.searchParams.set("description", book.description);
  if (book.release) url.searchParams.set("release", String(book.release));
  if (book.publisher) url.searchParams.set("publisher", book.publisher);
  if (book.thumbnail_link) url.searchParams.set("thumbnail_link", book.thumbnail_link);
  if (book.pages) url.searchParams.set("pages", String(book.pages));

  url.searchParams.set("ownership_status", String(book.ownership_status));
  if (book.owned_since && book.ownership_status == "owned") url.searchParams.set("owned_since", String(book.owned_since));
  url.searchParams.set("reading_status", String(book.reading_status));


  const res = await authState.authFetch(
    url.toString(),
    {
      method: "POST",
      body: JSON.stringify(book),
    },
    customFetch
  );
  if (!res.ok) {
    const errorData = await res.json().catch(() => ({}));
    throw new Error(errorData.error || `Failed to create book: ${res.statusText}`);
  }
  const data = await res.json();
  return data.data;
}

export async function removeBookById(id: number, customFetch?: typeof fetch): Promise<any> {
  const res = await authState.authFetch(`${API_BASE}/api/v1/books/${id}`, {method: "DELETE"}, customFetch);
  if (!res.ok) {
    throw new Error(`Failed to delete book: ${res.statusText}`);
  }
}

export async function lookupIsbn(isbn: string, customFetch?: typeof fetch): Promise<LookupBook> {
  const fetchFn = customFetch || fetch;
  const res = await fetchFn(`${API_BASE}/api/v1/lookup?isbn=${encodeURIComponent(isbn)}`);
  if (!res.ok) {
    const errorData = await res.json().catch(() => ({}));
    throw new Error(errorData.error || 'Book not found for this ISBN');
  }
  const data = await res.json();
  return data;
}

