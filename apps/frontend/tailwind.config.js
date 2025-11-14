/** @type {import('tailwindcss').Config} */
export default {
  content: ['./index.html', './src/**/*.{js,ts,jsx,tsx}'],
  theme: {
    extend: {
      colors: {
        border: 'hsl(214.3 31.8% 91.4%)',
        input: 'hsl(214.3 31.8% 91.4%)',
        ring: 'hsl(215.1 91.4% 52.9%)',

        background: 'hsl(0 0% 100%)',
        foreground: 'hsl(222.2 84% 4.9%)',

        primary: {
          DEFAULT: 'hsl(222.2 47.4% 11.2%)',
          foreground: 'hsl(210 40% 98%)',
        },
        secondary: {
          DEFAULT: 'hsl(210 40% 96.1%)',
          foreground: 'hsl(222.2 47.4% 11.2%)',
        },
      },
      borderRadius: {
        lg: `var(--radius)`,
        md: `calc(var(--radius) - 2px)`,
        sm: 'calc(var(--radius) - 4px)',
      },
      typography: ({ theme }) => ({
        DEFAULT: {
          css: {
            'code::before': {
              content: '""',
            },
            'code::after': {
              content: '""',
            },
            code: {
              backgroundColor: theme('colors.secondary.DEFAULT'),
              color: theme('colors.secondary.foreground'),
              fontWeight: '500',
              padding: '2px 6px',
              borderRadius: '6px',
              border: `1px solid ${theme('colors.border')}`,
            },
          },
        },
      }),
    },
  },
  plugins: [require('@tailwindcss/typography')],
};
