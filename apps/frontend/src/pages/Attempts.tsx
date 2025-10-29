import React from 'react';
import { useQuery } from '@tanstack/react-query';
import { fetchAttemptsForUser } from '../api/attempts';
import { Link } from 'react-router-dom';
import { Card } from '../components/ui';
import { CheckCircle2, XCircle, History } from 'lucide-react';

const AttemptStatus = ({ status }: { status: string }) => {
    if (status === 'ok') {
        return <span className="flex items-center gap-1.5 text-sm font-medium text-green-600"><CheckCircle2 className="w-4 h-4" /> Успешно</span>;
    }
    return <span className="flex items-center gap-1.5 text-sm font-medium text-red-600"><XCircle className="w-4 h-4" /> Ошибка</span>;
};


export default function Attempts() {
    const { data: attempts, isLoading } = useQuery({
        queryKey: ['myAttempts'],
        queryFn: fetchAttemptsForUser
    });

    if (isLoading) return <div>Загружаем историю...</div>;

    return (
        <div className="max-w-4xl mx-auto space-y-6">
            <h1 className="text-3xl font-bold tracking-tight">История попыток</h1>

            <Card>
                {attempts && attempts.length > 0 ? (
                    <div className="space-y-3">
                        {attempts.map((a) => (
                            <Link to={`/attempts/${a.id}`} key={a.id}>
                                <div className="p-4 bg-secondary rounded-lg flex items-center justify-between hover:bg-secondary/80 transition-colors">
                                    <div>
                                        <div className="font-semibold">Задача #{a.problemId}</div>
                                        <div className="text-sm text-gray-500 mt-1">
                                            {new Date(a.created_at ?? Date.now()).toLocaleString()}
                                        </div>
                                    </div>
                                    <div className="flex items-center gap-6">
                                        <span className="text-sm">{a.duration_ms} ms</span>
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
