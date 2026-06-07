import { Bell, Search, Star } from "lucide-react";
import { Link } from "react-router-dom";

export function Topbar() {
  return (
    <header className="sticky top-0 z-20 border-b border-border bg-background/95 py-4 backdrop-blur sm:px-6 lg:px-8">
      <section className="flex flex-col gap-3 p-4 md:flex-row md:items-center md:justify-between">
        <div className="relative w-full md:max-w-sm">
          <Search className="absolute left-3 top-1/2 size-4 -translate-y-1/2 text-muted-foreground" />
          <input
            aria-label="Search vocabulary"
            className="h-10 w-full rounded-xl border border-border bg-muted/40 pl-9 pr-4 text-sm outline-none transition focus:ring-2 focus:ring-primary/30"
            placeholder="Search vocabulary..."
            type="search"
          />
        </div>

        <div className="flex items-center gap-3">
          <div className="flex items-center gap-1.5 rounded-full border border-amber-200 bg-amber-50 px-3 py-1.5">
            <Star className="size-3.5 fill-amber-500 text-amber-500" />
            <span className="text-xs font-semibold text-amber-700">
              2,840 XP
            </span>
          </div>
          <button
            aria-label="Notifications"
            className="relative cursor-pointer rounded-xl p-2 transition-colors hover:bg-muted/60"
            type="button"
          >
            <Bell className="size-5 text-muted-foreground" />
            <span className="absolute right-1.5 top-1.5 size-2 rounded-full bg-red-500" />
          </button>
          <Link
            aria-label="Profile"
            className="flex size-9 cursor-pointer items-center justify-center rounded-full bg-indigo-100 text-sm font-semibold text-primary"
            to="/profile"
          >
            A
          </Link>
        </div>
      </section>
    </header>
  );
}
