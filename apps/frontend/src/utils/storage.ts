export function saveDraft(key: string, value: string) {
  try {
    localStorage.setItem(key, value);
  } catch (e) {
    console.error('Failed to save draft to localStorage', e);
  }
}

export function loadDraft(key: string): string | null {
  try {
    return localStorage.getItem(key);
  } catch (e) {
    console.error('Failed to load draft from localStorage', e);
    return null;
  }
}

export function removeDraft(key: string) {
  try {
    localStorage.removeItem(key);
  } catch (e) {
    console.error('Failed to remove draft from localStorage', e);
  }
}
