const typescript = require("@rollup/plugin-typescript");
const {uglify} = require("rollup-plugin-uglify");
const replace = require("@rollup/plugin-replace");

const plugins = [
    ...(process.env.NODE_ENV === 'production' ? [uglify()] : []),
];

module.exports = [
    {
        input: 'src/index.ts',
        output: [
            {
                file: 'dist/browserless/yopa.cjs.js',
                format: 'cjs',
                exports: 'auto',
                inlineDynamicImports: true,
                name: 'yopa'
            },
            {
                file: 'dist/browserless/yopa.esm.js',
                format: 'es',
                inlineDynamicImports: true,
                name: 'yopa'
            }
        ],
        plugins:  [
            typescript({
                tsconfig: 'tsconfig.json'
            }),
            ...plugins
        ],
        external: ['https']
    },
    {
        input: 'src/index.ts',
        output: [
            {
                file: 'dist/browser/yopa.umd.js',
                format: 'umd',
                inlineDynamicImports: true,
                name: 'yopa'
            },
            {
                file: 'dist/browser/yopa.js',
                format: 'iife',
                inlineDynamicImports: true,
                name: 'yopa'
            }
        ],
        plugins:  [
            typescript({
                tsconfig: 'tsconfig.browser.json'
            }),
            replace({BUILD_BROWSER: '"true"'}),
            ...plugins
        ]
    }
];