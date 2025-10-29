
import React, { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import api from '../api/client'
import { useAuth } from '../features/auth/hooks'

export default function Register() {
    const { login: doLogin } = useAuth()
    const navigate = useNavigate()

    const [loginVal, setLoginVal] = useState('')
    const [password, setPassword] = useState('')
    const [repeated, setRepeated] = useState('')
    const [loading, setLoading] = useState(false)
    const [error, setError] = useState<string | null>(null)

    const submit = async (e: React.FormEvent) => {
        e.preventDefault()
        setError(null)

        if (!loginVal || !password || !repeated) {
            setError('Заполни все поля')
            return
        }
        if (password !== repeated) {
            setError('Пароли не совпадают')
            return
        }

        try {
            setLoading(true)
            await api.post('/auth/register', { login: loginVal, password, repeated_password: repeated })
            await doLogin(loginVal, password)
            navigate('/', { replace: true })
        } catch (e: any) {
            setError(e?.response?.data?.message || e?.message || 'Не удалось зарегистрироваться')
        } finally {
            setLoading(false)
        }
    }

    return (
        <div className="min-h-screen flex items-center justify-center bg-gray-50">
            <div className="w-full max-w-md p-8 bg-white rounded-2xl shadow">
                <h1 className="text-2xl font-semibold mb-6">Регистрация</h1>
                <form onSubmit={submit} className="space-y-4">
                    <div>
                        <label className="block text-sm font-medium mb-1">Login</label>
                        <input
                            value={loginVal}
                            onChange={(e) => setLoginVal(e.target.value)}
                            className="w-full px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-brand"
                            placeholder="логин"
                        />
                    </div>

                    <div>
                        <label className="block text-sm font-medium mb-1">Password</label>
                        <input
                            value={password}
                            onChange={(e) => setPassword(e.target.value)}
                            type="password"
                            className="w-full px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-brand"
                            placeholder="пароль"
                        />
                    </div>

                    <div>
                        <label className="block text-sm font-medium mb-1">Repeat Password</label>
                        <input
                            value={repeated}
                            onChange={(e) => setRepeated(e.target.value)}
                            type="password"
                            className="w-full px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-brand"
                            placeholder="еще раз пароль"
                        />
                    </div>

                    {error && <div className="text-red-600 text-sm">{error}</div>}

                    <div className="flex items-center justify-between">
                        <button
                            type="submit"
                            disabled={loading}
                            className="px-4 py-2 rounded-lg bg-brand text-white font-medium disabled:opacity-60"
                        >
                            {loading ? 'Зарегистрироваться...' : 'Зарегистрироваться'}
                        </button>
                        <a href="/login" className="text-sm text-gray-600 hover:underline">
                            Уже есть аккаунт?
                        </a>
                    </div>
                </form>
            </div>
        </div>
    )
}
