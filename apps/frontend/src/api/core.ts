import type { Course, Lesson, Problem } from '../types';
import api from './client';

type CourseListResponse = {
  courses: Course[];
};

export async function fetchCourses(): Promise<Course[]> {
  const resp = await api.get<CourseListResponse>('/courses');
  return resp.data.courses;
}

export async function fetchCourse(id: number): Promise<Course | null> {
  const resp = await api.get<Course>(`/courses/${id}`);
  return resp.data;
}

export async function fetchLesson(courseId: number, lessonId: number): Promise<Lesson | null> {
  const resp = await api.get<Lesson>(`/courses/${courseId}/lessons/${lessonId}`);
  return resp.data;
}

export async function fetchProblem(
  courseId: number,
  problemId: number,
  language: string,
): Promise<Problem | null> {
  const resp = await api.get<Problem>(`/courses/${courseId}/problems/${problemId}`, {
    params: { language },
  });
  return resp.data;
}
