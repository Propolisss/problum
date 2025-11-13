import React from 'react';
import Editor from '@monaco-editor/react';
import type { editor } from 'monaco-editor';

import ruffInit, { Workspace, PositionEncoding } from '@astral-sh/ruff-wasm-web';
import initGoFmt, { format as formatGo } from '@wasm-fmt/gofmt/vite';

type Props = {
    value: string;
    language?: string;
    onChange?: (value: string) => void;
    height?: string;
    onMount?: (editor: editor.IStandaloneCodeEditor, monaco: any) => void;
    readOnly?: boolean;
};

let isGoFmtInitialized = false;

let ruffWorkspace: Workspace | null = null;

export default function CodeEditor({
    value,
    language = 'go',
    onChange,
    onMount,
    height = '400px',
    readOnly = false,
}: Props) {
    const handleEditorMount = (editorInstance: editor.IStandaloneCodeEditor, monacoInstance: any) => {
        if (!isGoFmtInitialized) {
            initGoFmt()
                .then(() => {
                    isGoFmtInitialized = true;
                })
                .catch(console.error);
        }

        if (!ruffWorkspace) {
            ruffInit()
                .then(() => {
                    const settings = Workspace.defaultSettings();
                    ruffWorkspace = new Workspace(settings, PositionEncoding.Utf16);
                    console.log('Ruff Workspace initialized successfully.');
                })
                .catch(console.error);
        }

        monacoInstance.languages.registerDocumentFormattingEditProvider('go', {
            async provideDocumentFormattingEdits(model: editor.ITextModel) {
                if (!isGoFmtInitialized) return [];
                const text = model.getValue();
                try {
                    const formatted = formatGo(text);
                    return [{ range: model.getFullModelRange(), text: formatted }];
                } catch (e) {
                    console.error('Go formatting error:', e);
                    return [];
                }
            },
        });

        monacoInstance.languages.registerDocumentFormattingEditProvider('python', {
            async provideDocumentFormattingEdits(model: editor.ITextModel) {
                if (!ruffWorkspace) {
                    console.warn('Ruff formatter is not ready yet.');
                    return [];
                }
                const text = model.getValue();
                try {
                    const formatted = ruffWorkspace.format(text);

                    return [{ range: model.getFullModelRange(), text: formatted }];
                } catch (e) {
                    console.error('Python formatting error:', e);
                    return [];
                }
            },
        });

        if (readOnly) {
            editorInstance.updateOptions({
                readOnly: true,
                domReadOnly: true,
                renderLineHighlight: 'none',
                scrollBeyondLastLine: false,
            });

            const dom = editorInstance.getDomNode();
            if (dom) {
                dom.setAttribute('tabindex', '-1');
            }

            const mouseDownDisposable = editorInstance.onMouseDown((e) => {
                e.event.preventDefault();
                e.event.stopPropagation();
            });

            const keyDownDisposable = editorInstance.onKeyDown((k) => {
                k.preventDefault();
            });

            const styleEl = document.createElement('style');
            styleEl.className = 'monaco-readonly-style';
            styleEl.textContent = `
        .monaco-editor .cursor { visibility: hidden !important; }
        .monaco-editor .view-lines { cursor: default !important; }
        .monaco-editor .selectionHighlight { display: none !important; }
      `;
            document.head.appendChild(styleEl);

            editorInstance.onDidDispose(() => {
                try {
                    mouseDownDisposable.dispose();
                } catch { }
                try {
                    keyDownDisposable.dispose();
                } catch { }
                try {
                    styleEl.remove();
                } catch { }
            });
        } else {
            editorInstance.updateOptions({
                readOnly: false,
            });
        }

        if (onMount) {
            onMount(editorInstance, monacoInstance);
        }
    };

    return (
        <div className="rounded-lg overflow-hidden border h-full">
            <Editor
                height={height}
                defaultLanguage={language}
                language={language}
                value={value}
                onChange={(v) => onChange?.(v ?? '')}
                onMount={handleEditorMount}
                options={{
                    minimap: { enabled: false },
                    fontSize: 14,
                    readOnly: readOnly,
                    domReadOnly: true,
                    renderControlCharacters: false,
                    lineNumbers: 'on',
                    scrollBeyondLastLine: false,
                }}
            />
        </div>
    );
}
