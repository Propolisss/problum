import React from 'react';
import { useParams } from 'react-router-dom';
import { useLesson } from '../features/courses/hooks';
import { Card } from '../components/ui';

export default function Lesson() {
    const { courseId, lessonId } = useParams();
    const { data: lesson, isLoading, isError } = useLesson(Number(courseId), Number(lessonId));

    if (isLoading) return <Card>Загружаем урок...</Card>;

    if (isError) {
        return (
            <Card className="text-center p-8">
                <h2 className="text-xl font-bold">Ошибка</h2>
                <p className="text-gray-600">Урок недоступен или не существует.</p>
            </Card>
        );
    }

    if (!lesson) return <Card>Урок не найден</Card>;

    return (
        <Card>
            <h1 className="text-3xl font-bold tracking-tight mb-4">{lesson.name}</h1>
            <div
                className="prose prose-lg max-w-none"
                dangerouslySetInnerHTML={{ __html: lesson.content }}
            />
        </Card>
    );
}
