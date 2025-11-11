import React from 'react';
import { useParams, Link } from 'react-router-dom';
import { useQuery } from '@tanstack/react-query';
import { fetchAttemptById } from '../api/attempts';
import { Card } from '../components/ui';
import { CheckCircle2, XCircle, ArrowLeft } from 'lucide-react';
import { formatMemory } from '../utils/formatters';

export default function AttemptDetail() {
    const { id } = useParams();
    const attemptId = Number(id);
    const { data, isLoading } = useQuery({
        queryKey: ['attempt', attemptId],
        queryFn: () => fetchAttemptById(attemptId),
        enabled: !!attemptId,
    });

    if (isLoading) return <div>Загружаем детали попытки...</div>;
    if (!data) return <div>Попытка не найдена</div>;

    const isSuccess = data.status === 'AC';
    const StatusIcon = isSuccess ? CheckCircle2 : (data.status === 'pending' ? ArrowLeft : XCircle);
    const statusColor = isSuccess ? 'text-green-600' : 'text-red-600';

    return (
        <div className="max-w-3xl mx-auto space-y-4">
            <Link to="/attempts" className="inline-flex items-center gap-2 text-sm text-primary font-medium hover:underline">
                <ArrowLeft className="w-4 h-4" />
                Назад к истории
            </Link>

            <h1 className="text-3xl font-bold tracking-tight">Попытка #{data.id}</h1>

            <Card>
                <div className="p-4 border-b">
                    <div className="flex items-center justify-between">
                        <div className="font-semibold text-lg">Задача #{data.problem_id}</div>
                        <div className={`flex items-center gap-2 font-bold ${statusColor}`}>
                            <StatusIcon className="w-5 h-5" />
                            <span>{data.status}</span>
                        </div>
                    </div>
                    <div className="text-sm text-gray-500 mt-2">
                        {new Date(data.created_at ?? Date.now()).toLocaleString()}
                    </div>
                </div>

                <div className="p-4 grid grid-cols-2 gap-4">
                    <div>
                        <div className="text-xs text-gray-500">Время выполнения</div>
                        <div className="text-sm font-medium">{(data.duration / 1_000_000).toFixed(2)} ms</div>
                    </div>
                    <div>
                        <div className="text-xs text-gray-500">Использовано памяти</div>
                        <div className="text-sm font-medium">{formatMemory(data.memory_usage)}</div>
                    </div>
                </div>

                {data.error_message && (
                    <div className="p-4 border-t">
                        <h3 className="font-semibold mb-2 text-red-600">Сообщение об ошибке</h3>
                        <pre className="bg-red-50 text-red-700 text-sm p-3 rounded-md whitespace-pre-wrap font-mono">
                            {data.error_message}
                        </pre>
                    </div>
                )}
            </Card>
        </div>
    );
}
