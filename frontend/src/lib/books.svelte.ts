import { fetchMyBooks, removeBookById, updateBook } from "./api";
import type { Book } from "./types";

class BookStore {
    books = $state<Book[]>([]);
    isLoaded = $state(false);
    isLoading = $state(false);

  // Fetch only if not already loaded
    async loadBooks(forceRefresh = false) {
        if (this.isLoaded && !forceRefresh) return;
        
        this.isLoading = true;
        try {
            this.books = await fetchMyBooks();
            this.isLoaded = true;
        } catch (err) {
            console.error('Failed to load books:', err);
        } finally {
            this.isLoading = false;
        }
    }

    async updateBook(book: Book) {
        const previousBook = this.books.filter((b) => b.id === book.id);
        this.books = this.books.map((b) => (b.id === book.id ? book : b));

        try {
            await updateBook(book);
        } catch (err) {
            console.error('Failed to sync update with backend, rolling back:', err);
            this.books = this.books.map((b) => (b.id === book.id ? previousBook[0] : b));
        }
    }

    async removeBook(book: Book) {
        this.books = this.books.filter((b) => (b.id !== book.id));

        try {
            await removeBookById(book.id);
        } catch (err) {
            console.error('Failed to sync update with backend, rolling back:', err);
        }
    }
}

export const bookStore = new BookStore();