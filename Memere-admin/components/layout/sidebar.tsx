"use client";

import { useState } from "react";
import { GraduationCap, PanelLeftClose, PanelLeftOpen } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Sheet, SheetContent, SheetTitle } from "@/components/ui/sheet";
import { SidebarNav } from "./sidebar-nav";
import { cn } from "@/lib/utils";

interface SidebarProps {
  appName?: string;
}

export function Sidebar({ appName = "Memere" }: SidebarProps) {
  const [collapsed, setCollapsed] = useState(false);

  return (
    <>
      {/* Desktop sidebar */}
      <aside
        className={cn(
          "hidden md:flex flex-col transition-all duration-200",
          collapsed ? "w-16" : "w-[220px]"
        )}
        style={{
          background: "hsl(var(--sidebar-bg))",
          color: "hsl(var(--sidebar-fg))",
        }}
      >
        {/* Logo / title */}
        <div
          className={cn(
            "flex h-16 items-center border-b px-3",
            collapsed ? "justify-center" : "justify-start"
          )}
          style={{ borderColor: "hsl(var(--sidebar-border))" }}
        >
          <div className="w-8 h-8 rounded-lg bg-foreground flex items-center justify-center shrink-0">
            <GraduationCap className="h-4 w-4 text-white" />
          </div>
          {!collapsed && (
            <span
              className="font-semibold text-sm ml-2.5 truncate"
              style={{ color: "hsl(var(--sidebar-fg))" }}
            >
              {appName}
            </span>
          )}
        </div>

        {/* Nav links */}
        <div className="flex-1 overflow-y-auto py-3">
          <SidebarNav collapsed={collapsed} />
        </div>

        {/* Divider */}
        <div style={{ borderColor: "hsl(var(--sidebar-border))", borderTopWidth: 1 }} />

        {/* Collapse toggle */}
        <div className={cn("p-3", collapsed ? "flex justify-center" : "flex justify-end")}>
          <Button
            variant="ghost"
            size="icon"
            aria-label={collapsed ? "Expand sidebar" : "Collapse sidebar"}
            onClick={() => setCollapsed((c) => !c)}
            className="hover:text-white"
            style={{ color: "hsl(var(--sidebar-muted))" }}
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
  appName = "Memere",
}: {
  open: boolean;
  onClose: () => void;
  appName?: string;
}) {
  return (
    <Sheet open={open} onOpenChange={(v) => !v && onClose()}>
      <SheetContent
        side="left"
        className="w-[220px] p-0 flex flex-col border-r-0"
        style={{
          background: "hsl(var(--sidebar-bg))",
          color: "hsl(var(--sidebar-fg))",
        }}
      >
        <SheetTitle
          className="flex h-16 items-center border-b px-3"
          style={{ borderColor: "hsl(var(--sidebar-border))" }}
        >
          <div className="w-8 h-8 rounded-lg bg-foreground flex items-center justify-center shrink-0">
            <GraduationCap className="h-4 w-4 text-white" />
          </div>
          <span
            className="font-semibold text-sm ml-2.5 truncate"
            style={{ color: "hsl(var(--sidebar-fg))" }}
          >
            {appName}
          </span>
        </SheetTitle>
        <div className="flex-1 overflow-y-auto py-3">
          <SidebarNav collapsed={false} />
        </div>
      </SheetContent>
    </Sheet>
  );
}
