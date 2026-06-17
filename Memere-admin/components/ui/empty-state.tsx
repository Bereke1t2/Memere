interface EmptyStateProps {
  title?: string;
  description?: string;
  children?: React.ReactNode;
}

export function EmptyState({
  title = "No results",
  description = "Nothing to show here yet.",
  children,
}: EmptyStateProps) {
  return (
    <div className="flex flex-col items-center justify-center gap-2 py-12 text-center">
      <p className="text-sm font-medium">{title}</p>
      <p className="text-sm text-muted-foreground">{description}</p>
      {children}
    </div>
  );
}
