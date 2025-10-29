import React from 'react'
import { NavLink } from 'react-router-dom'
import { useAuth } from '../../features/auth/hooks'
import { Button } from '../ui'

export default function Navbar() {
    const { logout } = useAuth()

    const linkClass = ({ isActive }: { isActive: boolean }) =>
        `text-sm font-medium transition-colors ${isActive ? 'text-brand-dark' : 'text-gray-500 hover:text-brand-dark'}`

    return (
        <header className="bg-white shadow-sm sticky top-0 z-10">
            <nav className="container mx-auto px-6 py-3 flex items-center justify-between">
                <NavLink to="/" className="text-xl font-bold text-brand-dark">
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

                <div>
                    <Button variant="ghost" onClick={logout} className="text-sm">
                        Выйти
                    </Button>
                </div>
            </nav>
        </header>
    )
}
