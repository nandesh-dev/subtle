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
            primary: '#48CB89',
            yellow: '#CBB148',
            orange: '#CB7A48',
            red: '#CB484F',
            gray: {
                830: '#D4D5D4',
                520: '#838683',
                190: '#2F3130',
                120: '#1E1F1E',
                80: '#191A19',
                50: '#121212',
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
