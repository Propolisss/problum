import React from 'react';
import { Outlet, useParams } from 'react-router-dom';
import { useCourse } from '../../features/courses/hooks';
import { Card } from '../ui';
import CourseSidebar from './CourseSidebar';

export default function CourseLayout() {
    const { courseId } = useParams();
    const { data: course, isLoading, isError } = useCourse(Number(courseId));

    if (isLoading) {
        return <div className="text-center p-8">Загрузка курса...</div>;
    }

    if (isError) {
        return (
            <div className="text-center p-8">
                <h2 className="text-xl font-bold">Ошибка</h2>
                <p className="text-gray-600">Курс недоступен или не существует.</p>
            </div>
        );
    }

    if (!course) {
        return <div className="text-center p-8">Не удалось загрузить курс.</div>;
    }

    if (!course.enrolled) {
        return (
            <div className="text-center p-8">
                <h2 className="text-xl font-bold">Доступ закрыт</h2>
                <p className="text-gray-600">Вы не записаны на этот курс. Вернитесь на главную страницу, чтобы записаться.</p>
            </div>
        )
    }

    return (
        <div className="grid grid-cols-1 lg:grid-cols-12 gap-8 max-w-screen-xl mx-auto">
            <div className="lg:col-span-3">
                <Card className="p-0 sticky top-24">
                    <CourseSidebar course={course} />
                </Card>
            </div>
            <div className="lg:col-span-9">
                <Outlet />
            </div>
        </div>
    );
}
