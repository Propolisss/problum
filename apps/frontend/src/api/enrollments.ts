import api from './client';

export async function enrollInCourse(courseId: number): Promise<void> {
  await api.post('/enrollments', { course_id: courseId });
}
