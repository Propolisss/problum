import React from 'react';
import { Moon, Sun } from 'lucide-react';
import { useTheme } from '../../features/theme/hooks';

export default function ThemeToggle() {
    const { theme, toggleTheme } = useTheme();
    const isDark = theme === 'dark';
    const Icon = isDark ? Sun : Moon;

    return (
        <button
            type="button"
            onClick={toggleTheme}
            className="
                inline-flex h-10 w-10 items-center justify-center rounded-md border border-border
                bg-card text-muted-foreground shadow-sm transition-colors
                hover:bg-secondary hover:text-foreground
                focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2
                focus-visible:ring-offset-background
            "
            aria-label={isDark ? 'Включить светлую тему' : 'Включить темную тему'}
            title={isDark ? 'Светлая тема' : 'Темная тема'}
        >
            <Icon className="h-4 w-4" />
        </button>
    );
}
