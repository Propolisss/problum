import React from 'react';
import { Navigate, useParams } from 'react-router-dom';
import { useCourse } from '../features/courses/hooks';

export default function CourseIndexPage() {
    const { courseId } = useParams();
    const { data: course, isLoading, isError } = useCourse(Number(courseId));

    if (isLoading) return <div>Загрузка...</div>;

    if (isError || !course || !course.lessons || course.lessons.length === 0) {
        return <Navigate to="/" replace />;
    }

    const sortedLessons = [...course.lessons].sort((a, b) => a.position - b.position);
    if (sortedLessons.length === 0) return <Navigate to="/" replace />;
    const firstLessonId = sortedLessons[0].id;

    return <Navigate to={`/courses/${courseId}/lessons/${firstLessonId}`} replace />;
}
