import React from 'react'
import type { APIAttempt } from '../api/attempts'
import { Card } from './ui'
import { CheckCircle2, XCircle, AlertTriangle, Loader2, Clock, MemoryStick } from 'lucide-react'
import { formatMemory } from '../utils/formatters'

type Props = {
    result: APIAttempt | null
    isPending: boolean
}

const Stat = ({ label, value, icon }: { label: string; value: string | number, icon: React.ReactNode }) => (
    <div className="flex items-center gap-2">
        {icon}
        <div>
            <div className="text-xs text-gray-500">{label}</div>
            <div className="text-sm font-medium text-foreground">{value}</div>
        </div>
    </div>
)

const statusConfig = {
    AC: { text: 'Решение принято', icon: CheckCircle2, color: 'text-green-600' },
    WA: { text: 'Неверный ответ', icon: XCircle, color: 'text-red-600' },
    CE: { text: 'Ошибка компиляции', icon: AlertTriangle, color: 'text-yellow-600' },
    RE: { text: 'Ошибка выполнения', icon: AlertTriangle, color: 'text-red-600' },
    TLE: { text: 'Превышен лимит времени', icon: Clock, color: 'text-red-600' },
    MLE: { text: 'Превышен лимит памяти', icon: MemoryStick, color: 'text-red-600' },
    pending: { text: 'В очереди...', icon: Loader2, color: 'text-gray-600' },
};

export default function SubmissionResult({ result, isPending }: Props) {
    if (isPending && !result) {
        return (
            <Card className="flex flex-col items-center justify-center min-h-[150px]">
                <Loader2 className="w-8 h-8 text-primary animate-spin" />
                <p className="mt-3 text-gray-600">Отправка решения...</p>
            </Card>
        )
    }

    if (!result) {
        return (
            <Card className="flex flex-col items-center justify-center min-h-[150px]">
                <p className="text-gray-600">Отправьте решение на проверку</p>
                <p className="text-xs text-gray-400 mt-2">Результаты появятся здесь</p>
            </Card>
        )
    }

    const currentStatus = statusConfig[result.status] || statusConfig.RE;
    const Icon = currentStatus.icon;

    return (
        <Card className="min-h-[150px]">
            <div className="flex items-center gap-3">
                <Icon className={`w-6 h-6 ${currentStatus.color} ${result.status === 'pending' || isPending ? 'animate-spin' : ''}`} />
                <h3 className={`font-semibold text-lg ${currentStatus.color}`}>
                    {currentStatus.text}
                </h3>
            </div>

            {result.status !== 'pending' && (
                <>
                    <div className="grid grid-cols-2 gap-4 mt-4 border-b pb-4 mb-4">
                        <Stat label="Время" value={`${(result.duration / 1_000_000).toFixed(2)} ms`} icon={<Clock className="w-4 h-4 text-gray-500" />} />
                        <Stat label="Память" value={formatMemory(result.memory_usage)} icon={<MemoryStick className="w-4 h-4 text-gray-500" />} />
                    </div>
                    {result.error_message && (
                        <div>
                            <h4 className="text-sm font-medium mb-2">Детали</h4>
                            <pre className="bg-secondary text-secondary-foreground text-xs p-3 rounded-md whitespace-pre-wrap font-mono">
                                {result.error_message}
                            </pre>
                        </div>
                    )}
                </>
            )}
        </Card>
    )
}
