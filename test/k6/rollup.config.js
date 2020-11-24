import commonjs from '@rollup/plugin-commonjs'
import json from '@rollup/plugin-json'
import resolve from '@rollup/plugin-node-resolve'
import babel from 'rollup-plugin-babel'
import { terser } from 'rollup-plugin-terser'
import multiInput from 'rollup-plugin-multi-input';
import path from 'path';

const extensions = ['.js', '.ts'];

export default [
  {
    input: ['src/test-*.ts'],
    external: id => ['k6'].some(m => id.includes(m)),
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
        transformOutputPath: (output, input) => path.basename(output),
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