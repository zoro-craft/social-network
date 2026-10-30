// Ported from frontend/src/utils/format.js
export default function formatDate(value: string | number): string {
  const date = new Date(value);
  if (isNaN(date.getTime())) return '';
  const diffMs = Date.now() - date.getTime();
  const minutes = Math.floor(diffMs / 60000);

  if (minutes < 1) return 'now';
  if (minutes < 60) return `${minutes}m`;
  const hours = Math.floor(minutes / 60);
  if (hours < 24) return `${hours}h`;
  const days = Math.floor(hours / 24);
  if (days < 7) return `${days}d`;
  const weeks = Math.floor(days / 7);
  if (weeks < 5) return `${weeks}w`;
  const months = Math.floor(days / 30);
  if (months < 12) return `${months}mo`;
  const years = Math.floor(days / 365);
  return `${years}y`;
}

// Matches the group event pill: "DD-MM-YYYY |HH:MM"
export function formatEventDateTime(dateStr?: string | null, timeStr?: string | null): string {
  if (!dateStr || !timeStr) return '';
  const [year, month, day] = dateStr.split('-');
  if (!year || !month || !day) return '';
  return `${day}-${month}-${year} |${timeStr}`;
}