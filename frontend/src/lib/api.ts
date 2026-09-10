import { authState, API_BASE } from './auth.svelte';
import type { Book, LookupBook, User } from './types';

// Helper to normalize books from backend to support both camelCase and snake_case properties
export function normalizeBook(raw: any): Book {
  return {
    ...raw,
    Isbn: raw.isbn || raw.Isbn || '',
    Title: raw.title || raw.Title || '',
    Author: raw.author || raw.Author || '',
    Description: raw.description || raw.Description || '',
    ThumbnailLink: raw.thumbnail_link || raw.ThumbnailLink || '',
    Pages: raw.pages || raw.Pages || 0,
    Release: raw.release || raw.Release || 0,
    Publisher: raw.publisher || raw.Publisher || '',
  };
}

export async function fetchBooks(customFetch?: typeof fetch): Promise<Book[]> {
  const res = await authState.authFetch(`${API_BASE}/api/v1/books`, {}, customFetch);
  if (!res.ok) {
    throw new Error(`Failed to fetch books: ${res.statusText}`);
  }
  const data = await res.json();
  const list = (data.data || []) as any[];
  return list.map(normalizeBook);
}

export async function fetchBookById(id: number, customFetch?: typeof fetch): Promise<Book> {
  const res = await authState.authFetch(`${API_BASE}/api/v1/books/${id}`, {}, customFetch);
  if (!res.ok) {
    throw new Error(`Failed to fetch book: ${res.statusText}`);
  }
  const data = await res.json();
  return normalizeBook(data.data);
}

export async function createBook(book: LookupBook, customFetch?: typeof fetch): Promise<Book> {
  const baseUrl = `${API_BASE}/api/v1/books`;
  const url = new URL(baseUrl, typeof window !== "undefined" ? window.location.origin : "http://localhost");
  if (book.isbn) url.searchParams.set("isbn", book.isbn);
  if (book.title) url.searchParams.set("title", book.title);
  if (book.author) url.searchParams.set("author", book.author);
  if (book.description) url.searchParams.set("description", book.description);
  if (book.thumbnail_link) url.searchParams.set("thumbnail_link", book.thumbnail_link);
  if (book.pages) url.searchParams.set("pages", String(book.pages));

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
  return normalizeBook(data.data);
}

export async function lookupIsbn(isbn: string, customFetch?: typeof fetch): Promise<LookupBook> {
  const fetchFn = customFetch || fetch;
  const res = await fetchFn(`${API_BASE}/api/v1/lookup?isbn=${encodeURIComponent(isbn)}`);
  if (!res.ok) {
    const errorData = await res.json().catch(() => ({}));
    throw new Error(errorData.error || 'Book not found for this ISBN');
  }
  const data = await res.json();
  return {
    isbn: data.Isbn || data.isbn || isbn,
    title: data.Title || data.title || '',
    author: data.Author || data.author || '',
    description: data.Description || data.description || '',
    thumbnail_link: data.ThumbnailLink || data.thumbnail_link || '',
    pages: data.Pages || data.pages || 0,
  };
}

