import commonjs from '@rollup/plugin-commonjs'
import json from '@rollup/plugin-json'
import resolve from '@rollup/plugin-node-resolve'
import babel from 'rollup-plugin-babel'
import { terser } from 'rollup-plugin-terser'
import multiInput from 'rollup-plugin-multi-input';
import utils from '@rollup/pluginutils';
import pkg from './package.json';

const extensions = ['.js', '.ts'];

export default [
  {
    input: ['src/test/**/*.ts'],
    external: utils.createFilter([
      'k6/**',
      ...Object.keys(pkg.devDependencies),
    ], null, { resolve: false }),
    output: [
      {
        dir: 'dist',
        format: 'cjs',
        exports: 'named',
        chunkFileNames: '_chunks/[name]-[hash].js'
      },
    ],
    plugins: [
      multiInput({
        transformOutputPath: (output, input) => `tests/${output.split('/').join('-')}`,
      }),
      json(),
      resolve(
        {
          extensions,
        }
      ),
      commonjs(),
      babel({
        extensions,
        include: ['src/**/*'],
      }),
      terser(),
    ],
  }
]