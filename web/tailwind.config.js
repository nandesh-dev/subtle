const defaultTheme = require('tailwindcss/defaultTheme')

/** @type {import('tailwindcss').Config} */
export default {
    content: ['./index.html', './src/**/*.{ts,tsx}'],
    theme: {
        extend: {
            fontFamily: {
                sans: ['"Sora"', ...defaultTheme.fontFamily.sans],
            },
        },
        colors: {
          primary: {
            1: "var(--primary-1-color)",
            2: "var(--primary-2-color)",
          },
          secondary: {
            1: "var(--secondary-1-color)",
            2: "var(--secondary-2-color)",
          },
          neutral: {
            1: "var(--neutral-1-color)",
            2: "var(--neutral-2-color)",
            3: "var(--neutral-3-color)",
          },
          text: {
            1: "var(--text-1-color)",
            2: "var(--text-2-color)",
          },
        },
        borderRadius: {
            xxs: '0.3rem',
            xs: '0.6rem',
            sm: '1.2rem',
            md: '1.8rem',
            lg: '2.4rem',
            xl: '3.6rem',
            xxl: '4.8rem',
        },
        spacing: {
            0: '0',
            xxs: '0.3rem',
            xs: '0.6rem',
            sm: '1.2rem',
            md: '1.8rem',
            lg: '2.4rem',
            xl: '3.6rem',
            xxl: '4.8rem',
        },
        fontSize: {
            xs: '1em',
            sm: '1.1em',
            md: '1.4em',
            lg: '1.7em',
            xl: '2.33em',
        },
    },
    plugins: [],
}
