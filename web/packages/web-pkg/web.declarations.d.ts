// This file must not export or import anything on top-level

declare module '*?worker' {
  const content: string
  export default content
}
