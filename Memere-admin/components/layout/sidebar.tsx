"use client";

import { useState } from "react";
import { PanelLeftClose, PanelLeftOpen } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Separator } from "@/components/ui/separator";
import { Sheet, SheetContent, SheetTitle } from "@/components/ui/sheet";
import { SidebarNav } from "./sidebar-nav";
import { cn } from "@/lib/utils";

interface SidebarProps {
  appName?: string;
}

export function Sidebar({ appName = "Memere Admin" }: SidebarProps) {
  const [collapsed, setCollapsed] = useState(false);

  return (
    <>
      {/* Desktop sidebar */}
      <aside
        className={cn(
          "hidden md:flex flex-col border-r bg-card transition-all duration-200",
          collapsed ? "w-14" : "w-56"
        )}
      >
        {/* Logo / title */}
        <div className={cn("flex h-14 items-center border-b px-3", collapsed && "justify-center")}>
          {!collapsed && (
            <span className="text-sm font-semibold truncate">{appName}</span>
          )}
        </div>

        {/* Nav links */}
        <div className="flex-1 overflow-y-auto py-3">
          <SidebarNav collapsed={collapsed} />
        </div>

        <Separator />

        {/* Collapse toggle */}
        <div className="flex justify-end p-2">
          <Button
            variant="ghost"
            size="icon"
            aria-label={collapsed ? "Expand sidebar" : "Collapse sidebar"}
            onClick={() => setCollapsed((c) => !c)}
          >
            {collapsed ? (
              <PanelLeftOpen className="h-4 w-4" />
            ) : (
              <PanelLeftClose className="h-4 w-4" />
            )}
          </Button>
        </div>
      </aside>
    </>
  );
}

/** Mobile drawer — opens from a hamburger in the Header. */
export function MobileSidebar({
  open,
  onClose,
  appName = "Memere Admin",
}: {
  open: boolean;
  onClose: () => void;
  appName?: string;
}) {
  return (
    <Sheet open={open} onOpenChange={(v) => !v && onClose()}>
      <SheetContent side="left" className="w-56 p-0 flex flex-col">
        <SheetTitle className="flex h-14 items-center border-b px-4 text-sm font-semibold">
          {appName}
        </SheetTitle>
        <div className="flex-1 overflow-y-auto py-3">
          <SidebarNav collapsed={false} />
        </div>
      </SheetContent>
    </Sheet>
  );
}
