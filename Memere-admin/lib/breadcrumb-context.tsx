"use client";

import { createContext, useContext, useState, useEffect } from "react";

interface BreadcrumbCtx {
  label: string;
  setLabel: (l: string) => void;
}

const BreadcrumbContext = createContext<BreadcrumbCtx>({
  label: "",
  setLabel: () => {},
});

export function BreadcrumbProvider({ children }: { children: React.ReactNode }) {
  const [label, setLabel] = useState("");
  return (
    <BreadcrumbContext.Provider value={{ label, setLabel }}>
      {children}
    </BreadcrumbContext.Provider>
  );
}

export function useBreadcrumb() {
  return useContext(BreadcrumbContext);
}

/** Drop inside any Server Component's output to set the header breadcrumb label. */
export function BreadcrumbSetter({ label }: { label: string }) {
  const { setLabel } = useBreadcrumb();
  useEffect(() => {
    setLabel(label);
    return () => setLabel("");
  }, [label, setLabel]);
  return null;
}
