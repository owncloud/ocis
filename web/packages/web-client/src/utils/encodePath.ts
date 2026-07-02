export const encodePath = (path: string) => {
  return encodeURIComponent(path).split('%2F').join('/')
}
