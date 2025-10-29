import React from 'react'

type Props = React.InputHTMLAttributes<HTMLInputElement>

export default function Input(props: Props) {
    return (
        <input
            {...props}
            className={`
        flex h-10 w-full rounded-md border border-input bg-background 
        px-3 py-2 text-sm 
        file:border-0 file:bg-transparent file:text-sm file:font-medium 
        placeholder:text-gray-400 
        focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring 
        disabled:cursor-not-allowed disabled:opacity-50
        ${props.className ?? ''}
      `}
        />
    )
}
