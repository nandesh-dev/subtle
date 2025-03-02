/**
 * @see https://prettier.io/docs/en/configuration.html
 * @type {import("prettier").Config}
 */
const config = {
    trailingComma: 'es5',
    tabWidth: 4,
    semi: false,
    singleQuote: true,
    plugins: [
        'prettier-plugin-tailwindcss',
        '@trivago/prettier-plugin-sort-imports',
    ],
    importOrder: [
        '<THIRD_PARTY_MODULES>',
        '^@/gen/(.*)$',
        '^@/src/sections/(.*)$',
        '^@/src/components/(.*)$',
        '^@/src/utility/(.*)$',
        '^[./]',
    ],
    importOrderSeparation: true,
    importOrderSortSpecifiers: true,
}

export default config
