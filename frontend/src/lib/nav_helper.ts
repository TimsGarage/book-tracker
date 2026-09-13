import { Library } from "lucide-svelte";
import { setContext } from "svelte";

export type NavState = {
  heading: string;
  showBackButton: boolean;
  onBack?: () => void;
  icon?: any;
  hideBottomNavigation?: boolean;
};

export let navStateDefaults: NavState = {
  heading: "",
  showBackButton: false,
  onBack: () => {},
  icon: Library,
  hideBottomNavigation: false,
}
