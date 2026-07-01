"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import {
  LayoutDashboard,
  Users,
  BookOpen,
  CreditCard,
  TrendingUp,
  Megaphone,
  DollarSign,
} from "lucide-react";
import { cn } from "@/lib/utils";
import { Tooltip, TooltipContent, TooltipTrigger } from "@/components/ui/tooltip";

const ADMIN_NAV = [
  { href: "/",              label: "Dashboard",     icon: LayoutDashboard },
  { href: "/users",         label: "Users",         icon: Users },
  { href: "/courses",       label: "Courses",       icon: BookOpen },
  { href: "/payments",      label: "Payments",      icon: CreditCard },
  { href: "/revenue",       label: "Revenue",       icon: TrendingUp },
  { href: "/announcements", label: "Announcements", icon: Megaphone },
];

const TEACHER_NAV = [
  { href: "/",            label: "Dashboard",  icon: LayoutDashboard },
  { href: "/my-courses",  label: "My Courses", icon: BookOpen },
  { href: "/earnings",    label: "Earnings",   icon: DollarSign },
];

interface SidebarNavProps {
  collapsed: boolean;
  role: string;
}

export function SidebarNav({ collapsed, role }: SidebarNavProps) {
  const pathname = usePathname();
  const items = role === "teacher" ? TEACHER_NAV : ADMIN_NAV;

  return (
    <nav aria-label="Main navigation" className="flex flex-col gap-1 px-2">
      {items.map(({ href, label, icon: Icon }) => {
        const isActive =
          href === "/" ? pathname === "/" : pathname.startsWith(href);

        const item = (
          <Link
            key={href}
            href={href}
            aria-current={isActive ? "page" : undefined}
            className={cn(
              "flex items-center gap-3 rounded-lg px-3 py-2 text-sm font-medium transition-all duration-150 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring",
              isActive
                ? "bg-sidebar-active-bg text-sidebar-active-fg shadow-sm"
                : "text-sidebar-muted hover:bg-sidebar-hover-bg hover:text-sidebar-fg",
              collapsed && "justify-center px-0 w-10 h-10 mx-auto"
            )}
          >
            <Icon className="h-4 w-4 shrink-0" aria-hidden="true" />
            {!collapsed && <span>{label}</span>}
          </Link>
        );

        if (collapsed) {
          return (
            <Tooltip key={href} delayDuration={0}>
              <TooltipTrigger asChild>{item}</TooltipTrigger>
              <TooltipContent side="right">{label}</TooltipContent>
            </Tooltip>
          );
        }

        return item;
      })}
    </nav>
  );
}
