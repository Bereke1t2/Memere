export const qk = {
  overview: (from: string, to: string) =>
    ["admin", "analytics", "overview", from, to] as const,

  engagement: () => ["admin", "analytics", "engagement"] as const,

  revenueBreakdown: (from: string, to: string) =>
    ["admin", "analytics", "revenue", from, to] as const,

  users: (params: { limit?: number; after?: string; role?: string } = {}) =>
    ["admin", "users", params] as const,

  user: (id: string) => ["admin", "users", id] as const,

  courses: (params: { limit?: number; after?: string } = {}) =>
    ["admin", "courses", params] as const,

  course: (id: string) => ["admin", "courses", id] as const,

  payments: (params: { limit?: number; after?: string; status?: string } = {}) =>
    ["admin", "payments", params] as const,

  payment: (id: string) => ["admin", "payments", id] as const,
} as const;
