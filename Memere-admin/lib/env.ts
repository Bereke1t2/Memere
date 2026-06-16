function requireEnv(key: string): string {
  const val = process.env[key];
  if (!val && process.env.NODE_ENV === "production") {
    throw new Error(`Missing required server environment variable: ${key}`);
  }
  return val ?? "";
}

export const env = {
  API_BASE_URL: requireEnv("API_BASE_URL"),
  COOKIE_SECRET: requireEnv("COOKIE_SECRET"),
} as const;
