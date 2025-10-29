import { useEffect, useRef, useState } from 'react';

export type AutosaveStatus = 'idle' | 'saving' | 'saved';

export function useAutosave(value: string, key: string, delay = 2000) {
  const timerRef = useRef<number | null>(null);
  const clearSavedTimerRef = useRef<number | null>(null);
  const [status, setStatus] = useState<AutosaveStatus>('idle');

  useEffect(() => {
    if (!key) return;

    setStatus('saving');

    if (timerRef.current) {
      window.clearTimeout(timerRef.current);
    }

    timerRef.current = window.setTimeout(() => {
      try {
        localStorage.setItem(key, value);
        setStatus('saved');
        if (clearSavedTimerRef.current) window.clearTimeout(clearSavedTimerRef.current);
        clearSavedTimerRef.current = window.setTimeout(() => setStatus('idle'), 1500);
      } catch (e) {
        setStatus('idle');
      }
    }, delay);

    return () => {
      if (timerRef.current) {
        window.clearTimeout(timerRef.current);
        timerRef.current = null;
      }
      if (clearSavedTimerRef.current) {
        window.clearTimeout(clearSavedTimerRef.current);
        clearSavedTimerRef.current = null;
      }
    };
  }, [value, key, delay]);

  return status;
}
