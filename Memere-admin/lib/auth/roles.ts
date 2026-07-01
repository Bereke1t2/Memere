import type { User } from "@/lib/api/schemas";

export function canManageContent(user: Pick<User, "role">): boolean {
  return user.role === "teacher" || user.role === "admin";
}
