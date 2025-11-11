export function formatMemory(bytes: number): string {
  if (bytes === 0) return '0 KB';

  const kb = bytes / 1024;
  if (kb < 1024) {
    return `${kb.toFixed(2)} KB`;
  }

  const mb = kb / 1024;
  return `${mb.toFixed(2)} MB`;
}
