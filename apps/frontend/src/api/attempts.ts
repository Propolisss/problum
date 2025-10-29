import api from './client';

export type APIAttempt = {
  id: number;
  user_id: number;
  problem_id: number;
  duration: number;
  memory_usage: number;
  language: string;
  code: string;
  status: 'pending' | 'AC' | 'WA' | 'CE' | 'RE' | 'TLE' | 'MLE';
  error_message: string | null;
  created_at: string;
  updated_at: string;
};

type SubmitResponse = {
  attempt_id: number;
};

export async function submitAttempt(
  courseId: number,
  problemId: number,
  language: string,
  code: string,
): Promise<SubmitResponse> {
  const resp = await api.post(`/courses/${courseId}/problems/${problemId}/submit`, {
    language,
    code,
  });
  return resp.data;
}

export async function fetchAttemptById(attemptId: number): Promise<APIAttempt> {
  const resp = await api.get(`/attempts/${attemptId}`);
  return resp.data;
}

export async function fetchAttemptsForProblem(
  courseId: number,
  problemId: number,
): Promise<APIAttempt[]> {
  const resp = await api.get(`/courses/${courseId}/problems/${problemId}/attempts`);
  return (resp.data as { attempts: APIAttempt[] }).attempts;
}

export async function fetchAttemptsForUser(): Promise<APIAttempt[]> {
  const resp = await api.get('/attempts');
  return (resp.data as { attempts: APIAttempt[] }).attempts;
}
