"use client";

import { useState } from "react";
import { Sidebar, MobileSidebar } from "./sidebar";
import { Header } from "./header";
import type { User } from "@/lib/api/schemas";

interface DashboardShellProps {
  user: User;
  children: React.ReactNode;
}

export function DashboardShell({ user, children }: DashboardShellProps) {
  const [mobileOpen, setMobileOpen] = useState(false);

  return (
    <div className="flex h-screen overflow-hidden bg-background">
      <Sidebar role={user.role} />
      <MobileSidebar open={mobileOpen} onClose={() => setMobileOpen(false)} role={user.role} />
      <div className="flex flex-1 flex-col min-w-0">
        <Header user={user} onMenuClick={() => setMobileOpen(true)} />
        <main
          id="main-content"
          className="flex-1 overflow-auto bg-[linear-gradient(180deg,hsl(var(--background)),hsl(var(--muted)/0.36))]"
        >
          <div className="mx-auto w-full max-w-screen-2xl px-4 py-5 sm:px-6 lg:px-8 lg:py-7">
            {children}
          </div>
        </main>
      </div>
    </div>
  );
}
