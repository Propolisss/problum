import React from 'react';
import { Link } from 'react-router-dom';
import { useCourses, useEnrollInCourse } from '../features/courses/hooks';
import { Button, Card } from '../components/ui';

export default function Dashboard() {
    const { data: courses, isLoading, isError } = useCourses();
    const enrollMutation = useEnrollInCourse();

    const handleEnrollClick = (courseId: number) => {
        enrollMutation.mutate(courseId);
    };

    if (isLoading) return <div>Загружаем курсы...</div>;
    if (isError) return <div>Ошибка загрузки курсов</div>;

    return (
        <div className="space-y-6">
            <h1 className="text-3xl font-bold tracking-tight">Доступные курсы</h1>
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                {courses?.map((course) => (
                    <Card key={course.id} className="flex flex-col">
                        <h2 className="text-xl font-semibold text-foreground">
                            {course.enrolled ? (
                                <Link to={`/courses/${course.id}`} className="hover:underline">
                                    {course.name}
                                </Link>
                            ) : (
                                <span>{course.name}</span>
                            )}
                        </h2>
                        <p className="text-sm text-gray-600 mt-2 flex-grow">{course.description}</p>

                        <div className="flex items-center justify-between mt-6">
                            <div className="flex flex-wrap gap-2">
                                {course.tags.map(tag => (
                                    <span key={tag} className="text-xs px-2 py-1 bg-secondary rounded-md text-secondary-foreground">
                                        {tag}
                                    </span>
                                ))}
                            </div>

                            {course.enrolled ? (
                                <Link to={`/courses/${course.id}`} className="text-sm text-primary font-semibold hover:underline shrink-0 ml-4">
                                    Продолжить
                                </Link>
                            ) : (
                                <Button
                                    variant="secondary"
                                    className="text-sm shrink-0 ml-4"
                                    onClick={() => handleEnrollClick(course.id)}
                                    disabled={enrollMutation.isPending && enrollMutation.variables === course.id}
                                >
                                    {enrollMutation.isPending && enrollMutation.variables === course.id ? 'Запись...' : 'Записаться'}
                                </Button>
                            )}
                        </div>
                    </Card>
                ))}
            </div>
        </div>
    );
}
