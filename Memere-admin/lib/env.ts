function requireEnv(key: string): string {
  const val = process.env[key];
  if (!val && process.env.NODE_ENV === "production") {
    throw new Error(`Missing required server environment variable: ${key}`);
  }
  return val ?? "";
}

// Lazy getters: validation runs at request time, not at module evaluation / build.
export const env = {
  get API_BASE_URL() {
    return requireEnv("API_BASE_URL");
  },
  get COOKIE_SECRET() {
    return requireEnv("COOKIE_SECRET");
  },
};
