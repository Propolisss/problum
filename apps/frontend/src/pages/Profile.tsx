import React from 'react';
import { useQuery } from '@tanstack/react-query';
import { fetchAttemptsForUser } from '../api/attempts';
import { Card, Button } from '../components/ui';
import { useAuth } from '../features/auth/hooks';
import { Link } from 'react-router-dom';
import { CheckCircle2, XCircle, History } from 'lucide-react';

const AttemptStatus = ({ status }: { status: string }) => {
    if (status === 'ok') {
        return <span className="flex items-center gap-1.5 text-sm font-medium text-green-600"><CheckCircle2 className="w-4 h-4" /> Успешно</span>;
    }
    return <span className="flex items-center gap-1.5 text-sm font-medium text-red-600"><XCircle className="w-4 h-4" /> Ошибка</span>;
};

export default function Profile() {
    const { logout } = useAuth();
    const { data: attempts, isLoading, isError } = useQuery({
        queryKey: ['myAttempts'],
        queryFn: fetchAttemptsForUser
    });

    const successfulAttempts = attempts?.filter(a => a.status === 'ok').length ?? 0;
    const recentAttempts = attempts?.slice(0, 5) ?? [];

    if (isLoading) return <div>Загружаем профиль...</div>;
    if (isError) return <div>Ошибка загрузки</div>;

    return (
        <div className="max-w-4xl mx-auto space-y-6">
            <h1 className="text-3xl font-bold tracking-tight">Мой профиль</h1>

            <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
                <Card className="md:col-span-1 flex flex-col items-center text-center">
                    <div className="w-24 h-24 rounded-full bg-secondary flex items-center justify-center mb-4">
                        <span className="text-4xl font-bold text-secondary-foreground">U</span>
                    </div>
                    <div className="font-semibold text-lg">Пользователь</div>
                    <p className="text-sm text-gray-500">На платформе с {new Date().toLocaleDateString()}</p>
                    <Button variant="ghost" onClick={logout} className="mt-4 text-red-600">
                        Выйти из аккаунта
                    </Button>
                </Card>

                <Card className="md:col-span-2">
                    <h2 className="font-semibold text-lg mb-4">Статистика</h2>
                    <div className="grid grid-cols-2 gap-4">
                        <div className="p-4 bg-secondary rounded-lg">
                            <div className="text-sm text-gray-600">Всего попыток</div>
                            <div className="text-3xl font-bold">{attempts?.length ?? 0}</div>
                        </div>
                        <div className="p-4 bg-secondary rounded-lg">
                            <div className="text-sm text-gray-600">Успешных решений</div>
                            <div className="text-3xl font-bold text-green-600">{successfulAttempts}</div>
                        </div>
                    </div>
                </Card>
            </div>

            <Card>
                <div className="flex items-center justify-between mb-4">
                    <h2 className="font-semibold text-lg">Последняя активность</h2>
                    <Link to="/attempts" className="text-sm font-medium text-primary hover:underline">
                        Вся история
                    </Link>
                </div>

                {recentAttempts.length > 0 ? (
                    <div className="space-y-3">
                        {recentAttempts.map((a) => (
                            <Link to={`/attempts/${a.id}`} key={a.id}>
                                <div className="p-3 bg-secondary rounded-md flex items-center justify-between hover:bg-secondary/80 transition-colors">
                                    <div>
                                        <div className="font-medium">Задача #{a.problemId}</div>
                                        <div className="text-xs text-gray-500 mt-1">{new Date(a.created_at ?? Date.now()).toLocaleString()}</div>
                                    </div>
                                    <AttemptStatus status={a.status} />
                                </div>
                            </Link>
                        ))}
                    </div>
                ) : (
                    <div className="text-center py-8 text-gray-500">
                        <History className="w-12 h-12 mx-auto mb-2" />
                        <p>Вы еще не решили ни одной задачи.</p>
                        <p className="text-sm">Самое время начать!</p>
                    </div>
                )}
            </Card>
        </div>
    );
}
