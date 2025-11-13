import type { APIAttempt } from '../api/attempts';

export type ActivityData = {
  heatmapValues: { date: Date; count: number }[];
  totalActiveDays: number;
  maxStreak: number;
  startDate: Date;
  endDate: Date;
};

function toYMDLocal(date: Date): string {
  const y = date.getFullYear();
  const m = date.getMonth() + 1;
  const d = date.getDate();
  const mm = m < 10 ? `0${m}` : `${m}`;
  const dd = d < 10 ? `0${d}` : `${d}`;
  return `${y}-${mm}-${dd}`;
}

function localDateFromYMD(ymd: string): Date {
  const [y, m, d] = ymd.split('-').map(Number);
  return new Date(y, m - 1, d);
}

export function processActivityData(attempts: APIAttempt[]): ActivityData {
  if (!attempts || attempts.length === 0) {
    const endDate = new Date();
    const startDate = new Date();
    startDate.setFullYear(endDate.getFullYear() - 1);
    return {
      heatmapValues: [],
      totalActiveDays: 0,
      maxStreak: 0,
      startDate,
      endDate,
    };
  }

  const countsByDate = new Map<string, number>();
  for (const attempt of attempts) {
    const dt = new Date(attempt.created_at);
    const key = toYMDLocal(dt);
    countsByDate.set(key, (countsByDate.get(key) || 0) + 1);
  }

  const totalActiveDays = countsByDate.size;

  let maxStreak = 0;
  if (totalActiveDays > 0) {
    const sortedKeys = Array.from(countsByDate.keys()).sort();
    let currentStreak = 1;
    maxStreak = 1;

    const msPerDay = 24 * 60 * 60 * 1000;
    for (let i = 1; i < sortedKeys.length; i++) {
      const currDate = localDateFromYMD(sortedKeys[i]);
      const prevDate = localDateFromYMD(sortedKeys[i - 1]);

      const diffDays = Math.round((currDate.getTime() - prevDate.getTime()) / msPerDay);

      if (diffDays === 1) {
        currentStreak++;
      } else {
        currentStreak = 1;
      }
      if (currentStreak > maxStreak) maxStreak = currentStreak;
    }
  }

  const heatmapValues = Array.from(countsByDate.entries()).map(([ymd, count]) => ({
    date: localDateFromYMD(ymd),
    count,
  }));

  const endDate = new Date();
  const startDate = new Date();
  startDate.setFullYear(endDate.getFullYear() - 1);

  return {
    heatmapValues,
    totalActiveDays,
    maxStreak,
    startDate,
    endDate,
  };
}
