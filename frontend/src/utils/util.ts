import { goto } from "$app/navigation";

export function handleBack() {
    if (window.history.length > 1) {
        window.history.back();
    } else {
        goto('/');
    }
}