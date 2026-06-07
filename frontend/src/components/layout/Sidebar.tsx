import {
  BarChart3,
  BookOpen,
  // Brain,
  Star,
  Target,
  Trophy,
  // Zap,
} from "lucide-react";
import { NavLink, useLocation } from "react-router-dom";

import { cn } from "@/lib/utils";

const navigationItems = [
  {
    label: "Dashboard",
    to: "/dashboard",
    icon: BarChart3,
    match: (pathname: string) => pathname === "/dashboard",
  },
  {
    label: "Roadmap",
    to: "/roadmap",
    icon: Target,
    match: (pathname: string) =>
      pathname === "/roadmap" || pathname.startsWith("/topics/"),
  },
  // {
  //   label: "Lessons",
  //   to: "/roadmap",
  //   icon: BookOpen,
  //   match: (pathname: string) =>
  //     pathname.startsWith("/lessons/") &&
  //     !pathname.endsWith("/flashcards") &&
  //     !pathname.endsWith("/quiz"),
  // },
  // {
  //   label: 'Flashcards',
  //   to: '/reviews',
  //   icon: Brain,
  //   match: (pathname: string) =>
  //     pathname === '/reviews' || pathname.endsWith('/flashcards'),
  // },
  // {
  //   label: 'Quiz',
  //   to: '/lessons/1/quiz',
  //   icon: Zap,
  //   match: (pathname: string) => pathname.endsWith('/quiz'),
  // },
  {
    label: "Profile",
    to: "/profile",
    icon: Trophy,
    match: (pathname: string) => pathname === "/profile",
  },
];

function getNavLinkClass(isActive: boolean) {
  return cn(
    "flex w-full cursor-pointer items-center gap-3 rounded-xl px-3 py-2.5 text-sm transition-colors",
    isActive
      ? "bg-accent font-medium text-primary"
      : "text-muted-foreground hover:bg-muted/60 hover:text-foreground",
  );
}

export function Sidebar() {
  const { pathname } = useLocation();

  return (
    <aside className="fixed left-0 top-0 z-30 hidden min-h-screen w-60 flex-col border-r border-border bg-white md:flex">
      <div className="flex h-16 items-center border-b border-border px-5">
        <NavLink className="flex cursor-pointer items-center gap-2" to="/">
          <span className="flex size-8 items-center justify-center rounded-lg bg-primary text-white">
            <BookOpen className="size-[15px]" aria-hidden="true" />
          </span>
          <span className="font-bold text-foreground">LexPath</span>
        </NavLink>
      </div>

      <nav className="flex-1 space-y-1 px-3 py-4" aria-label="Dashboard">
        {navigationItems.map((item) => (
          <NavLink
            className={() => getNavLinkClass(item.match(pathname))}
            key={item.to}
            to={item.to}
          >
            <item.icon className="size-[18px]" aria-hidden="true" />
            {item.label}
          </NavLink>
        ))}
      </nav>

      <div className="mx-3 mb-4 rounded-xl border border-amber-200 bg-amber-50 p-3">
        <div className="mb-1 flex items-center gap-2">
          <Star className="size-3.5 fill-amber-500 text-amber-500" />
          <span className="text-xs font-semibold text-amber-800">
            Daily Challenge
          </span>
        </div>
        <p className="mb-2 text-xs text-amber-700">
          Learn 10 words before midnight!
        </p>
        <div className="h-1.5 overflow-hidden rounded-full bg-amber-200">
          <div className="h-full w-3/5 rounded-full bg-amber-500" />
        </div>
        <div className="mt-1 text-xs text-amber-700">6 / 10 words</div>
      </div>
    </aside>
  );
}

export function DashboardMobileNav() {
  const { pathname } = useLocation();

  return (
    <nav className="border-b border-border bg-card px-4 py-3 md:hidden">
      <div className="flex gap-2 overflow-x-auto">
        {navigationItems.map((item) => (
          <NavLink
            className={() =>
              cn(getNavLinkClass(item.match(pathname)), "shrink-0 gap-2")
            }
            key={item.to}
            to={item.to}
          >
            <item.icon className="size-4" aria-hidden="true" />
            {item.label}
          </NavLink>
        ))}
      </div>
    </nav>
  );
}
