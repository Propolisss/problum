import React from 'react'
import { Outlet } from 'react-router-dom'
import Navbar from './Navbar'

export default function Layout() {
    return (
        <div className="min-h-screen bg-gray-100">
            <Navbar />
            <main>
                <div className="container mx-auto py-8 px-6">
                    <Outlet />
                </div>
            </main>
        </div>
    )
}
