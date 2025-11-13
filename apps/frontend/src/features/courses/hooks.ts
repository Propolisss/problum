import { useQuery, useMutation, useQueryClient, type UseQueryOptions } from '@tanstack/react-query';
import { fetchCourses, fetchCourse, fetchLesson, fetchProblem } from '../../api/core';
import { enrollInCourse } from '../../api/enrollments';
import { fetchAttemptsForProblem } from '../../api/attempts';
import axios from 'axios';

function useApiQuery<T>(options: UseQueryOptions<T, Error, T, any[]>) {
  return useQuery<T, Error, T, any[]>({
    ...options,
    retry: (failureCount, error) => {
      if (axios.isAxiosError(error) && error.response?.status === 403) {
        return false;
      }
      return failureCount < 3;
    },
  });
}

export function useCourses() {
  return useQuery({
    queryKey: ['courses'],
    queryFn: fetchCourses,
    staleTime: 1000 * 60,
  });
}

export function useCourse(id: number) {
  return useApiQuery({
    queryKey: ['course', id],
    queryFn: () => fetchCourse(id),
    staleTime: 1000 * 60,
    enabled: !!id,
  });
}

export function useLesson(courseId: number, lessonId: number) {
  return useApiQuery({
    queryKey: ['course', courseId, 'lesson', lessonId],
    queryFn: () => fetchLesson(courseId, lessonId),
    enabled: !!courseId && !!lessonId,
  });
}

export function useProblem(courseId: number, problemId: number, language: string) {
  return useApiQuery({
    queryKey: ['course', courseId, 'problem', problemId, language],
    queryFn: () => fetchProblem(courseId, problemId, language),
    enabled: !!courseId && !!problemId && !!language,
  });
}

export function useEnrollInCourse() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (courseId: number) => enrollInCourse(courseId),
    onSuccess: (data, courseId) => {
      queryClient.invalidateQueries({ queryKey: ['courses'] });
      queryClient.invalidateQueries({ queryKey: ['course', courseId] });
    },
  });
}

export function useProblemAttempts(courseId: number, problemId: number) {
  return useQuery({
    queryKey: ['attempts', 'problem', problemId],
    queryFn: () => fetchAttemptsForProblem(courseId, problemId),
    enabled: !!courseId && !!problemId,
  });
}
