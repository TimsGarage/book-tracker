import { goto } from "$app/navigation";
import { BookCheck, BookDown, Bookmark, BookMarked, BookmarkX, BookOpenText, Library } from "lucide-svelte";
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

export const read_options: {
  title: string,
  value: BookReadingStatus
  icon: any;
}[] = [
    {
      title: "Ungelesen",
      value: "unread",
      icon: BookMarked,
    },
    {
      title: "Am lesen",
      value: "reading",
      icon: BookOpenText,
    },
    {
      title: "Gelesen",
      value: "read",
      icon: BookCheck,
    },
  ];

export const owning_options: {
  title: string;
  value: BookOwnershipStatus
  icon: any
}[] = [
    {
      title: "Owned",
      value: "owned",
      icon: Library,
    },
    {
      title: "Unowned",
      value: "unowned",
      icon: BookmarkX,
    },
    {
      title: "Wishlist",
      value: "wishlist",
      icon: Bookmark,
    },
    {
      title: "Ausgeliehen",
      value: "borrowed",
      icon: BookDown,
    },
  ];