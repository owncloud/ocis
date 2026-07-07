// This file must not export or import anything on top-level

declare module '*.vue' {
  import { defineComponent } from 'vue'
  const component: ReturnType<typeof defineComponent>
  export default component
}

declare module 'join-path' {
  const joinPath: (...args: string[]) => string
  export default joinPath
}

declare module 'qs' {
  export function parse(str: string): Record<string, string | string[]>
  export function stringify(obj: Record<string, any>, opts?: Record<string, any>): string
  export function parse(str: string, opts?: Record<string, any>): Record<string, any>
}
