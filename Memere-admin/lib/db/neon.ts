import "server-only";
import { Pool } from "pg";

const connectionString =
  process.env.DATABASE_URL ||
  "postgresql://neondb_owner:npg_ki4GCJ8cOhPe@ep-blue-unit-b2ryzoq0-pooler.c-6.eu-central-1.aws.neon.tech/neondb?sslmode=require";

declare global {
  // eslint-disable-next-line no-var
  var __neon_pool: Pool | undefined;
}

export const dbPool: Pool =
  globalThis.__neon_pool ??
  new Pool({
    connectionString,
    ssl: true,
    max: 10,
    idleTimeoutMillis: 30000,
    connectionTimeoutMillis: 10000,
  });

if (process.env.NODE_ENV !== "production") {
  globalThis.__neon_pool = dbPool;
}
