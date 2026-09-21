import { goto } from "$app/navigation";
import { BookCheck, BookDown, Bookmark, BookMarked, BookmarkX, BookOpenText, Library, Star } from "lucide-svelte";
import type { BookOwnershipStatus, BookReadingStatus } from "./types";

export function handleBack() {
    if (window.history.length > 1) {
        window.history.back();
    } else {
        goto('/');
    }
}


export function getTodayString(): string {
  const today = new Date();
  const year = today.getFullYear();
  const month = String(today.getMonth() + 1).padStart(2, "0");
  const day = String(today.getDate()).padStart(2, "0");
  return `${year}-${month}-${day}`;
}

export const UNREAD_ICON = BookMarked;
export const READING_ICON = BookOpenText;
export const READ_ICON = BookCheck;

export const read_options: {
  title: string,
  value: BookReadingStatus
  icon: any;
}[] = [
    {
      title: "Ungelesen",
      value: "unread",
      icon: UNREAD_ICON,
    },
    {
      title: "Am lesen",
      value: "reading",
      icon: READING_ICON,
    },
    {
      title: "Gelesen",
      value: "read",
      icon: READ_ICON,
    },
  ];


export const OWNED_ICON = Library;
export const UNOWNED_ICON = Bookmark;
export const WISHLIST_ICON = Star;
export const BORROWED_ICON = BookDown;

export const owning_options: {
  title: string;
  value: BookOwnershipStatus
  icon: any
}[] = [
    {
      title: "Owned",
      value: "owned",
      icon: OWNED_ICON,
    },
    {
      title: "Unowned",
      value: "unowned",
      icon: UNOWNED_ICON,
    },
    {
      title: "Wishlist",
      value: "wishlist",
      icon: WISHLIST_ICON,
    },
    {
      title: "Borrowed",
      value: "borrowed",
      icon: BORROWED_ICON,
    },
  ];