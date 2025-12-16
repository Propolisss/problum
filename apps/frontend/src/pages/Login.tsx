import React, { useState } from 'react'
import { useNavigate, useLocation, Link } from 'react-router-dom'
import { useAuth } from '../features/auth/hooks'
import { Button, Input } from '../components/ui'

export default function Login() {
    const { login } = useAuth()
    const navigate = useNavigate()
    const location = useLocation() as any
    const from = location.state?.from?.pathname || '/'

    const [loginVal, setLoginVal] = useState('')
    const [password, setPassword] = useState('')
    const [loading, setLoading] = useState(false)
    const [error, setError] = useState<string | null>(null)

    const submit = async (e: React.FormEvent) => {
        e.preventDefault()
        setError(null)

        if (!loginVal || !password) {
            setError('Заполните все поля')
            return
        }

        try {
            setLoading(true)
            await login(loginVal, password)
            navigate(from, { replace: true })
        } catch (e: any) {
            setError('Произошла ошибка. Попробуйте позже.')
            console.error('Login error:', e)
        } finally {
            setLoading(false)
        }
    }

    return (
        <div className="min-h-screen flex items-center justify-center bg-gray-50">
            <div className="w-full max-w-md p-8 bg-white rounded-2xl shadow">
                <h1 className="text-2xl font-semibold mb-6">Вход</h1>
                <form onSubmit={submit} className="space-y-4">
                    <div>
                        <label className="block text-sm font-medium mb-1">Login</label>
                        <Input
                            value={loginVal}
                            onChange={(e) => setLoginVal(e.target.value)}
                            placeholder="логин"
                        />
                    </div>

                    <div>
                        <label className="block text-sm font-medium mb-1">Password</label>
                        <Input
                            value={password}
                            onChange={(e) => setPassword(e.target.value)}
                            type="password"
                            placeholder="пароль"
                        />
                    </div>

                    {error && <div className="text-red-600 text-sm">{error}</div>}

                    <div className="pt-2 space-y-3">
                        <Button
                            type="submit"
                            disabled={loading}
                            className="w-full"
                        >
                            {loading ? 'Вхожу...' : 'Войти'}
                        </Button>
                        <div className="text-center">
                            <Link to="/register" className="text-sm text-gray-600 hover:underline">
                                Регистрация
                            </Link>
                        </div>
                    </div>
                </form>
            </div>
        </div>
    )
}
