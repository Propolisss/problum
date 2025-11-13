import React, { useMemo, useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { fetchAttemptsForUser } from '../api/attempts';
import { Card, Button } from '../components/ui';
import { useAuth } from '../features/auth/hooks';
import { Link } from 'react-router-dom';
import { History, CheckCircle2, XCircle, Clock, MemoryStick } from 'lucide-react';
import { useProfile } from '../features/user/hooks';
import { processActivityData } from '../utils/activityProcessor';
import ActivityHeatmap from '../components/ActivityHeatmap';
import 'react-tooltip/dist/react-tooltip.css';
import { formatMemory } from '../utils/formatters';

const AttemptStatus = ({ status }: { status: string }) => {
    if (status === 'AC') {
        return <span className="flex items-center gap-1.5 text-sm font-medium text-green-600"><CheckCircle2 className="w-4 h-4" /> Успешно</span>;
    }
    return <span className="flex items-center gap-1.5 text-sm font-medium text-red-600"><XCircle className="w-4 h-4" /> Ошибка</span>;
};

export default function Profile() {
    const { logout } = useAuth();
    const { data: profile, isLoading: isLoadingProfile } = useProfile();
    const { data: attempts, isLoading: isLoadingAttempts } = useQuery({
        queryKey: ['myAttempts'],
        queryFn: fetchAttemptsForUser,
    });

    const [showAllAttempts, setShowAllAttempts] = useState(false);

    const activityData = useMemo(() => {
        return attempts ? processActivityData(attempts) : null;
    }, [attempts]);

    const sortedAttempts = useMemo(() => {
        if (!attempts) return [];
        return [...attempts].sort((a, b) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime());
    }, [attempts]);

    const visibleAttempts = showAllAttempts ? sortedAttempts : sortedAttempts.slice(0, 5);

    const successfulAttempts = attempts?.filter(a => a.status === 'AC').length ?? 0;

    if (isLoadingProfile || isLoadingAttempts) {
        return (
            <div className="flex items-center justify-center h-64">
                <div className="text-lg text-gray-600">Загружаем профиль...</div>
            </div>
        );
    }

    return (
        <div className="max-w-4xl mx-auto space-y-6">
            <div className="flex items-center justify-between">
                <h1 className="text-3xl font-bold tracking-tight">Мой профиль</h1>
                <Button variant="ghost" onClick={logout} className="text-red-600">
                    Выйти
                </Button>
            </div>

            <Card>
                <div className="flex items-center gap-4">
                    <div className="w-16 h-16 rounded-full bg-secondary flex items-center justify-center">
                        <span className="text-3xl font-bold text-secondary-foreground">
                            {profile?.login.charAt(0).toUpperCase()}
                        </span>
                    </div>
                    <div>
                        <div className="font-semibold text-xl">{profile?.login}</div>
                        <p className="text-sm text-gray-500">
                            На платформе с {profile ? new Date(profile.created_at).toLocaleDateString() : '...'}
                        </p>
                    </div>
                </div>
            </Card>

            <Card>
                <h2 className="font-semibold text-lg mb-4">Статистика решений</h2>
                <div className="grid grid-cols-2 gap-4">
                    <div className="p-4 bg-secondary rounded-lg">
                        <div className="text-sm text-gray-600">Всего попыток</div>
                        <div className="text-3xl font-bold">{attempts?.length ?? 0}</div>
                    </div>
                    <div className="p-4 bg-secondary rounded-lg">
                        <div className="text-sm text-gray-600">Успешных</div>
                        <div className="text-3xl font-bold text-green-600">{successfulAttempts}</div>
                    </div>
                </div>
            </Card>

            <Card>
                <h2 className="font-semibold text-lg mb-4">Активность за последний год</h2>
                {activityData ? (
                    <ActivityHeatmap data={activityData} />
                ) : (
                    <div className="text-center py-8 text-gray-500">
                        <History className="w-12 h-12 mx-auto mb-2" />
                        <p>Нет данных об активности.</p>
                    </div>
                )}
            </Card>

            <Card>
                <h2 className="font-semibold text-lg mb-4">Недавние попытки</h2>
                {visibleAttempts && visibleAttempts.length > 0 ? (
                    <div className="space-y-3">
                        {visibleAttempts.map((a) => (
                            <Link to={`/attempts/${a.id}`} key={a.id}>
                                <div className="p-4 bg-secondary rounded-lg flex items-center justify-between hover:bg-secondary/80 transition-colors">
                                    <div>
                                        <div className="font-semibold">Задача #{a.problem_id}</div>
                                        <div className="text-sm text-gray-500 mt-1">
                                            {new Date(a.created_at).toLocaleString()}
                                        </div>
                                    </div>
                                    <div className="hidden sm:flex items-center gap-6">
                                        <div className="flex items-center gap-1.5 text-sm text-gray-600" title="Время выполнения">
                                            <Clock className="w-4 h-4" />
                                            <span>{(a.duration / 1_000_000).toFixed(2)} ms</span>
                                        </div>
                                        <div className="flex items-center gap-1.5 text-sm text-gray-600" title="Потребление памяти">
                                            <MemoryStick className="w-4 h-4" />
                                            <span>{formatMemory(a.memory_usage)}</span>
                                        </div>
                                    </div>
                                    <AttemptStatus status={a.status} />
                                </div>
                            </Link>
                        ))}
                        {sortedAttempts.length > 5 && (
                            <div className="pt-2 text-center">
                                <Button variant="secondary" onClick={() => setShowAllAttempts(!showAllAttempts)}>
                                    {showAllAttempts ? 'Скрыть' : `Показать еще ${sortedAttempts.length - 5}`}
                                </Button>
                            </div>
                        )}
                    </div>
                ) : (
                    <div className="text-center py-12 text-gray-500">
                        <History className="w-16 h-16 mx-auto mb-3" />
                        <h3 className="font-semibold">История пуста</h3>
                        <p className="text-sm mt-1">Как только вы начнете решать задачи, ваши попытки появятся здесь.</p>
                    </div>
                )}
            </Card>
        </div>
    );
}
