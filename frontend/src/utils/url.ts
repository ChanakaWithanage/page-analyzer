export function isLikelyUrl(s: string): boolean {
  const trimmed = s.trim();
  if (!trimmed || trimmed.length > 2048) return false;
  // allow missing scheme locally; server will add https:// when needed
  const candidate = trimmed.includes("://") ? trimmed : `https://${trimmed}`;
  try {
    const u = new URL(candidate);
    return u.protocol === "http:" || u.protocol === "https:";
  } catch {
    return false;
  }
}
