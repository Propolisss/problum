import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { fetchCourses, fetchCourse, fetchLesson, fetchProblem } from '../../api/core';
import { enrollInCourse } from '../../api/enrollments';
import { fetchAttemptsForProblem } from '../../api/attempts';

export function useCourses() {
  return useQuery({
    queryKey: ['courses'],
    queryFn: fetchCourses,
    staleTime: 1000 * 60,
  });
}

export function useCourse(id: number) {
  return useQuery({
    queryKey: ['course', id],
    queryFn: () => fetchCourse(id),
    staleTime: 1000 * 60,
  });
}

export function useLesson(courseId: number, lessonId: number) {
  return useQuery({
    queryKey: ['course', courseId, 'lesson', lessonId],
    queryFn: () => fetchLesson(courseId, lessonId),
    enabled: !!courseId && !!lessonId,
  });
}

export function useProblem(courseId: number, problemId: number, language: string) {
  return useQuery({
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
