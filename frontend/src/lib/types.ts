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

export type BookOwnershipStatus = "owned" | "unowned" | "wishlist" | "borrowed";
export type BookReadingStatus = "unread" | "reading" | "read";

export interface Book {
  id: number;
  isbn: string;
  title: string;
  author: string;
  description?: string;
  release?: string;
  publisher?: string;
  thumbnail_link?: string;
  pages?: number;
  
  ownership_status: BookOwnershipStatus;
  owned_since?: string;
  reading_status: BookReadingStatus;

  user_id?: number;
  created_at?: string;
  updated_at?: string;
}

export interface LookupBook {
  isbn: string;
  title: string;
  author: string;
  description?: string;
  release?: string;
  publisher?: string;
  thumbnail_link?: string;
  pages?: number;
  
  owned?: boolean;
  borrowed?: boolean;
  wishlist?: boolean;

  read?: boolean;
  reading?: boolean;
}

export interface ApiError {
  error: string;
}
