import React from 'react'

export default function Card({ children, className = '' }: { children: React.ReactNode; className?: string }) {
    return (
        <div
            className={`
        p-6 bg-card text-card-foreground 
        border border-border rounded-lg shadow-sm 
        ${className}
      `}
        >
            {children}
        </div>
    )
}
