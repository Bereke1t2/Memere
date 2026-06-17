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
    <div className="flex h-screen overflow-hidden">
      <Sidebar />
      <MobileSidebar open={mobileOpen} onClose={() => setMobileOpen(false)} />
      <div className="flex flex-1 flex-col min-w-0">
        <Header user={user} onMenuClick={() => setMobileOpen(true)} />
        <main id="main-content" className="flex-1 overflow-auto p-6">
          {children}
        </main>
      </div>
    </div>
  );
}
