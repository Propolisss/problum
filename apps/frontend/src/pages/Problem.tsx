import { useEffect, useMemo, useState, useRef } from 'react';
import { Link, useParams } from 'react-router-dom';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { useProblem, useProblemAttempts, useCourse } from '../features/courses/hooks';
import { submitAttempt, fetchAttemptById, type APIAttempt } from '../api/attempts';
import { loadDraft, saveDraft } from '../utils/storage';
import CodeEditor from '../components/CodeEditor';
import SubmissionResult from '../components/SubmissionResult';
import { Card, Button, Select } from '../components/ui';
import { Clock, MemoryStick, BrainCircuit, History, CheckCircle2, XCircle, Sparkles, Loader2 } from 'lucide-react';
import type { editor } from 'monaco-editor';

const AttemptStatusIcon = ({ status }: { status: string }) => {
    if (status === 'AC') {
        return <CheckCircle2 className="w-4 h-4 text-green-600" />;
    }
    if (status === 'pending') {
        return <Loader2 className="w-4 h-4 text-gray-500 animate-spin" />;
    }
    return <XCircle className="w-4 h-4 text-red-600" />;
};

const supportedLanguages = [
    { value: 'go', label: 'Go' },
    { value: 'python', label: 'Python' },
];

export default function Problem() {
    const { courseId, problemId } = useParams();
    const queryClient = useQueryClient();
    const numericCourseId = Number(courseId);
    const numericProblemId = Number(problemId);

    const { data: course } = useCourse(numericCourseId);

    const problemMetadata = useMemo(() => {
        if (!course) return null;
        for (const lesson of course.lessons) {
            const p = lesson.problems?.find(p => p.id === numericProblemId);
            if (p) return p;
        }
        return null;
    }, [course, numericProblemId]);

    const availableLanguages = problemMetadata?.languages || [];

    const [language, setLanguage] = useState<string>('');

    useEffect(() => {
        if (availableLanguages.length > 0 && !language) {
            setLanguage(availableLanguages[0]);
        }
    }, [availableLanguages, language]);

    const { data: problem, isLoading: isLoadingProblem, isRefetching, isError } = useProblem(numericCourseId, numericProblemId, language);
    const { data: pastAttempts } = useProblemAttempts(numericCourseId, numericProblemId);

    const storageKey = useMemo(() => `draft_problem_${numericProblemId}_${language}`, [numericProblemId, language]);
    const [code, setCode] = useState<string>('');
    const [currentAttemptId, setCurrentAttemptId] = useState<number | null>(null);

    const editorRef = useRef<editor.IStandaloneCodeEditor | null>(null);
    const isCodeInitialized = useRef(false);

    useEffect(() => {
        isCodeInitialized.current = false;
        setCode('// Загрузка шаблона...');
    }, [language]);

    useEffect(() => {
        if (problem?.template?.code && !isCodeInitialized.current) {
            const draft = loadDraft(storageKey);
            setCode(draft ?? problem.template.code);
            isCodeInitialized.current = true;
            setTimeout(() => editorRef.current?.getAction('editor.action.formatDocument')?.run(), 200);
        }
    }, [problem, storageKey]);

    useEffect(() => {
        if (!isCodeInitialized.current) return;
        const handler = setTimeout(() => {
            if (code !== problem?.template?.code) {
                saveDraft(storageKey, code);
            }
        }, 1000);
        return () => clearTimeout(handler);
    }, [code, storageKey, problem?.template?.code]);

    const handleEditorDidMount = (editorInstance: editor.IStandaloneCodeEditor) => {
        editorRef.current = editorInstance;
        if (isCodeInitialized.current) {
            setTimeout(() => editorRef.current?.getAction('editor.action.formatDocument')?.run(), 200);
        }
    };

    const handleFormatClick = () => editorRef.current?.getAction('editor.action.formatDocument')?.run();

    const submissionMutation = useMutation({
        mutationFn: (vars: { language: string; code: string }) =>
            submitAttempt(numericCourseId, numericProblemId, vars.language, vars.code),
        onSuccess: (data) => {
            setCurrentAttemptId(data.attempt_id);
            queryClient.invalidateQueries({ queryKey: ['attempts', 'problem', numericProblemId] });
        },
    });

    const { data: attemptResult, isFetching: isPolling } = useQuery({
        queryKey: ['attempt', currentAttemptId],
        queryFn: () => fetchAttemptById(currentAttemptId!),
        enabled: !!currentAttemptId,
        refetchInterval: (query) => {
            const data = query.state.data;
            return data?.status === 'pending' ? 1000 : false;
        },
        refetchOnWindowFocus: false,
    });

    useEffect(() => {
        if (attemptResult && attemptResult.status !== 'pending') {
            queryClient.invalidateQueries({ queryKey: ['attempts', 'problem', numericProblemId] });
        }
    }, [attemptResult, queryClient, numericProblemId]);


    const onSubmit = () => {
        setCurrentAttemptId(null);
        submissionMutation.mutate({ language, code });
    };

    if (isLoadingProblem || !problemMetadata) return <Card>Загружаем задачу...</Card>;
    if (isError) {
        return (
            <Card className="text-center p-8">
                <h2 className="text-xl font-bold">Ошибка</h2>
                <p className="text-gray-600">Задача недоступна или не существует.</p>
            </Card>
        );
    }

    if (!problem) return <Card>Не удалось загрузить задачу</Card>;

    const isSubmitting = submissionMutation.isPending || (isPolling && attemptResult?.status === 'pending');

    return (
        <div className="grid grid-cols-1 xl:grid-cols-5 gap-6">
            <div className="xl:col-span-2 space-y-4">
                <h1 className="text-3xl font-bold tracking-tight">{problem.name}</h1>
                <div className="flex items-center gap-6 text-sm text-gray-600">
                    <span className="flex items-center gap-1.5"><BrainCircuit className="w-4 h-4" /> <span className="capitalize">{problem.difficulty}</span></span>
                    <span className="flex items-center gap-1.5"><Clock className="w-4 h-4" /> {problem.time_limit / 1_000_000_000}s</span>
                    <span className="flex items-center gap-1.5"><MemoryStick className="w-4 h-4" /> {problem.memory_limit / 1024 / 1024} MB</span>
                </div>
                <Card>
                    <div
                        className="prose max-w-none prose-pre:bg-secondary prose-pre:p-4 prose-pre:rounded-lg"
                        dangerouslySetInnerHTML={{ __html: problem.statement }}
                    />
                </Card>
                <Card>
                    <h3 className="font-semibold mb-3 flex items-center gap-2"><History className="w-5 h-5" />Прошлые попытки</h3>
                    {pastAttempts && pastAttempts.length > 0 ? (
                        <div className="space-y-2">
                            {[...pastAttempts].sort((a, b) => b.id - a.id).slice(0, 5).map(att => (
                                <Link key={att.id} to={`/attempts/${att.id}`}>
                                    <div className="p-2 bg-secondary rounded-md flex justify-between items-center text-sm hover:bg-secondary/80">
                                        <span className="font-mono text-xs">#{att.id} - {new Date(att.created_at).toLocaleString()}</span>
                                        <AttemptStatusIcon status={att.status} />
                                    </div>
                                </Link>
                            ))}
                        </div>
                    ) : (
                        <p className="text-sm text-gray-500">Вы еще не решали эту задачу.</p>
                    )}
                </Card>
            </div>

            <div className="xl:col-span-3 space-y-4 sticky top-24">
                <Card className="flex flex-col h-[70vh]">
                    <div className="flex items-center justify-between pb-4 border-b mb-4">
                        <Select
                            value={language}
                            onChange={(e) => setLanguage(e.target.value)}
                            disabled={availableLanguages.length <= 1}
                        >
                            {availableLanguages.map(lang => (
                                <option key={lang} value={lang}>
                                    {supportedLanguages.find(l => l.value === lang)?.label || lang}
                                </option>
                            ))}
                        </Select>
                        <Button variant="ghost" onClick={handleFormatClick} className="flex items-center gap-2">
                            <Sparkles className="w-4 h-4" />
                            Форматировать
                        </Button>
                    </div>

                    <div className="flex-grow min-h-0">
                        <CodeEditor
                            value={code}
                            language={language}
                            onChange={setCode}
                            height="100%"
                            onMount={handleEditorDidMount}
                        />
                    </div>
                    <div className="flex justify-end pt-4 border-t mt-4">
                        <Button onClick={onSubmit} disabled={isSubmitting || isRefetching}>
                            {isSubmitting ? 'Проверяется...' : (isRefetching ? 'Загрузка...' : 'Отправить')}
                        </Button>
                    </div>
                </Card>
                <div>
                    <SubmissionResult
                        isPending={isSubmitting}
                        result={attemptResult ?? null}
                    />
                </div>
            </div>
        </div>
    );
}
