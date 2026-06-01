import React from 'react'
import { NavLink } from 'react-router-dom'
import { useAuth } from '../../features/auth/hooks'
import { Button } from '../ui'
import ThemeToggle from '../theme/ThemeToggle'

export default function Navbar() {
    const { logout } = useAuth()

    const linkClass = ({ isActive }: { isActive: boolean }) =>
        `text-sm font-medium transition-colors ${isActive ? 'text-primary' : 'text-muted-foreground hover:text-foreground'}`

    return (
        <header className="sticky top-0 z-10 border-b border-border bg-card/90 shadow-sm backdrop-blur">
            <nav className="container mx-auto px-6 py-3 flex items-center justify-between">
                <NavLink to="/" className="text-xl font-bold text-foreground">
                    Problum
                </NavLink>

                <div className="flex items-center gap-6">
                    <NavLink to="/" className={linkClass} end>
                        Курсы
                    </NavLink>
                    <NavLink to="/attempts" className={linkClass}>
                        Мои Попытки
                    </NavLink>
                    <NavLink to="/profile" className={linkClass}>
                        Профиль
                    </NavLink>
                </div>

                <div className="flex items-center gap-2">
                    <ThemeToggle />
                    <Button variant="ghost" onClick={logout} className="text-sm">
                        Выйти
                    </Button>
                </div>
            </nav>
        </header>
    )
}
