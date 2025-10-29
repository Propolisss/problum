import React from 'react';
import { NavLink, useParams } from 'react-router-dom';
import type { Course, Problem } from '../../types';
import { BookText, Code2 } from 'lucide-react';

type Props = {
    course: Course;
};

export default function CourseSidebar({ course }: Props) {
    const { lessonId, problemId } = useParams();

    const baseLinkClasses = "flex items-center gap-3 px-3 py-2 rounded-md text-sm font-medium transition-colors";
    const inactiveClasses = "text-gray-600 hover:bg-secondary";
    const activeClasses = "bg-secondary text-primary font-semibold";

    const shouldExpandProblems = (lesson: { id: number; problems?: Problem[] }): boolean => {
        if (String(lesson.id) === lessonId) {
            return true;
        }
        if (problemId && lesson.problems?.some(p => String(p.id) === problemId)) {
            return true;
        }
        return false;
    };

    return (
        <aside className="w-full">
            <div className="p-4 border-b">
                <h2 className="font-bold text-lg text-foreground">{course.name}</h2>
                <p className="text-xs text-gray-500">Оглавление курса</p>
            </div>
            <nav className="p-4 space-y-1">
                {course.lessons?.map((lesson) => (
                    <div key={lesson.id}>
                        <NavLink
                            to={`/courses/${course.id}/lessons/${lesson.id}`}
                            className={({ isActive }) => `${baseLinkClasses} ${isActive ? activeClasses : inactiveClasses}`}
                        >
                            <BookText className="w-4 h-4" />
                            <span>{lesson.name}</span>
                        </NavLink>

                        {shouldExpandProblems(lesson) && lesson.problems && lesson.problems.length > 0 && (
                            <div className="pl-6 mt-1 space-y-1 border-l-2 ml-4 py-2">
                                {lesson.problems.map(problem => (
                                    <NavLink
                                        key={problem.id}
                                        to={`/courses/${course.id}/problems/${problem.id}`}
                                        className={({ isActive }) => `${baseLinkClasses} ${isActive ? activeClasses : inactiveClasses}`}
                                    >
                                        <Code2 className="w-4 h-4" />
                                        <span>{problem.name}</span>
                                    </NavLink>
                                ))}
                            </div>
                        )}
                    </div>
                ))}
            </nav>
        </aside>
    );
}
