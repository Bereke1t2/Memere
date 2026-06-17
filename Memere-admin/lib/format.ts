/**
 * The analytics/overview endpoints (/admin/analytics/*) do not return a currency
 * field in v1. ETB (Ethiopian Birr) is the platform's operating currency and is
 * used as a constant for those endpoints only. All other money values (payments,
 * course prices) read currency directly from the data field.
 */
export const PLATFORM_CURRENCY = "ETB";

export function formatMoney(value: string, currency: string): string {
  const num = parseFloat(value);
  if (isNaN(num)) return `${currency} 0.00`;
  try {
    return new Intl.NumberFormat("en-US", {
      style: "currency",
      currency,
      minimumFractionDigits: 2,
    }).format(num);
  } catch {
    return `${currency} ${num.toFixed(2)}`;
  }
}

export function formatPercent(value: number): string {
  return `${(value * 100).toFixed(1)}%`;
}

export function formatCount(value: number): string {
  return new Intl.NumberFormat("en-US").format(value);
}

export function formatDate(iso: string | null | undefined): string {
  if (!iso) return "—";
  return new Intl.DateTimeFormat("en-US", {
    year: "numeric",
    month: "short",
    day: "numeric",
  }).format(new Date(iso));
}
