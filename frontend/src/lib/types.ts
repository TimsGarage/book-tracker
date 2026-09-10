export interface User {
  id: number;
  username: string;
}

export interface AuthResponse {
  token: string;
  user: User;
}

export interface LoginCredentials {
  username: string;
  password: string;
}

export interface Book {
  id?: number;
  isbn: string;
  title: string;
  author: string;
  description?: string;
  thumbnail_link?: string;
  pages?: number;
  owned?: boolean;
  read?: boolean;
  user_id?: number;
  created_at?: string;
  updated_at?: string;

  // Compatibility aliases for existing Svelte components
  Isbn?: string;
  Title?: string;
  Author?: string;
  Publisher?: string;
  Release?: number;
  Description?: string;
  ThumbnailLink?: string;
  Pages?: number;
}

export interface LookupBook {
  isbn: string;
  title: string;
  author: string;
  description?: string;
  thumbnail_link?: string;
  pages?: number;
}

export interface ApiError {
  error: string;
}

