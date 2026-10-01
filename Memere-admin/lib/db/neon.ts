import "server-only";
import { Pool } from "pg";

const connectionString =
  process.env.DATABASE_URL ||
  "postgresql://neondb_owner:npg_ki4GCJ8cOhPe@ep-blue-unit-b2ryzoq0-pooler.c-6.eu-central-1.aws.neon.tech/neondb?sslmode=require";

// Global pool instance cached across hot-reloads and serverless invocations
declare global {
  // eslint-disable-next-line no-var
  var __neon_pool: Pool | undefined;
}

export const dbPool: Pool =
  globalThis.__neon_pool ??
  new Pool({
    connectionString,
    ssl: { rejectUnauthorized: false },
    max: 10,
    idleTimeoutMillis: 30000,
    connectionTimeoutMillis: 10000,
  });

if (process.env.NODE_ENV !== "production") {
  globalThis.__neon_pool = dbPool;
}
