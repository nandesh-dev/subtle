import type { Config } from 'tailwindcss'

export default {
    content: ['./index.html', './src/**/*.{ts,tsx}'],
    theme: {
        extend: {
            animation: {
                loader: 'loader 3s infinite',
            },
            keyframes: {
                loader: {
                    '0%': {
                        backgroundPosition: '-150% 0, -150% 0',
                    },
                    '66%': {
                        backgroundPosition: '251% 0, -150% 0',
                    },
                    '100%': {
                        backgroundPosition: '251% 0, 251% 0',
                    },
                },
            },
            fontFamily: {
                sans: [
                    '"Sora"',
                    'ui-sans-serif',
                    'system-ui',
                    'sans-serif',
                    '"Apple Color Emoji"',
                    '"Segoe UI Emoji"',
                    '"Segoe UI Symbol"',
                    '"Noto Color Emoji"',
                ],
            },
        },
        colors: {
            primary: {
                1: 'var(--primary-1-color)',
                2: 'var(--primary-2-color)',
            },
            secondary: {
                1: 'var(--secondary-1-color)',
                2: 'var(--secondary-2-color)',
            },
            neutral: {
                1: 'var(--neutral-1-color)',
                2: 'var(--neutral-2-color)',
                3: 'var(--neutral-3-color)',
            },
            text: {
                1: 'var(--text-1-color)',
                2: 'var(--text-2-color)',
            },
        },
        borderRadius: {
            xs: '0.3rem',
            sm: '0.6rem',
            md: '1.1rem',
            lg: '1.6rem',
            xl: '2.2rem',
            '2xl': '4.4rem',
        },
        spacing: {
            0: '0',
            xs: '0.3rem',
            sm: '0.6rem',
            md: '1.1rem',
            lg: '1.6rem',
            xl: '2.2rem',
            '2xl': '4.4rem',
        },
        fontSize: {
            xs: '0.9em',
            sm: '1em',
            default: '1em',
            md: '1.1em',
            lg: '1.6em',
            xl: '2em',
            '2xl': '4em',
        },
    },
    plugins: [],
} satisfies Config
