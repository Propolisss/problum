import React, { useMemo } from 'react';
import { useQuery } from '@tanstack/react-query';
import { fetchAttemptsForUser } from '../api/attempts';
import { Link } from 'react-router-dom';
import { Card } from '../components/ui';
import { CheckCircle2, XCircle, History, MemoryStick, Clock } from 'lucide-react';
import { formatMemory } from '../utils/formatters';

const AttemptStatus = ({ status }: { status: string }) => {
    if (status === 'AC') {
        return <span className="flex items-center gap-1.5 text-sm font-medium text-green-600"><CheckCircle2 className="w-4 h-4" /> Успешно</span>;
    }
    return <span className="flex items-center gap-1.5 text-sm font-medium text-red-600"><XCircle className="w-4 h-4" /> Ошибка</span>;
};


export default function Attempts() {
    const { data: attempts, isLoading } = useQuery({
        queryKey: ['myAttempts'],
        queryFn: fetchAttemptsForUser
    });

    const sortedAttempts = useMemo(() => {
        if (!attempts) return [];
        return [...attempts].sort((a, b) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime());
    }, [attempts]);

    if (isLoading) return <div>Загружаем историю...</div>;

    return (
        <div className="max-w-4xl mx-auto space-y-6">
            <h1 className="text-3xl font-bold tracking-tight">История попыток</h1>

            <Card>
                {sortedAttempts && sortedAttempts.length > 0 ? (
                    <div className="space-y-3">
                        {sortedAttempts.map((a) => (
                            <Link to={`/attempts/${a.id}`} key={a.id}>
                                <div className="p-4 bg-secondary rounded-lg flex items-center justify-between hover:bg-secondary/80 transition-colors">
                                    <div>
                                        <div className="font-semibold">Задача #{a.problem_id}</div>
                                        <div className="text-sm text-gray-500 mt-1">
                                            {new Date(a.created_at ?? Date.now()).toLocaleString()}
                                        </div>
                                    </div>
                                    <div className="flex items-center gap-6">
                                        <div className="flex items-center gap-1.5 text-sm text-gray-600" title="Время выполнения">
                                            <Clock className="w-4 h-4" />
                                            <span>{(a.duration / 1_000_000).toFixed(2)} ms</span>
                                        </div>
                                        <div className="flex items-center gap-1.5 text-sm text-gray-600" title="Потребление памяти">
                                            <MemoryStick className="w-4 h-4" />
                                            <span>{formatMemory(a.memory_usage)}</span>
                                        </div>
                                        <AttemptStatus status={a.status} />
                                    </div>
                                </div>
                            </Link>
                        ))}
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
