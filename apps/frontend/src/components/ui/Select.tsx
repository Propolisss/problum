import React from 'react';

type Props = React.SelectHTMLAttributes<HTMLSelectElement> & {
    children: React.ReactNode;
};

export default function Select({ children, className = '', ...rest }: Props) {
    return (
        <select
            className={`
        h-9 rounded-md border border-input bg-transparent px-3 py-1 text-sm shadow-sm 
        focus:outline-none focus:ring-1 focus:ring-ring
        disabled:cursor-not-allowed disabled:opacity-50
        ${className}
      `}
            {...rest}
        >
            {children}
        </select>
    );
}
