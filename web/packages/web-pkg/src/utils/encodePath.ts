export const encodePath = (path = ''): string => {
  return path.split('/').map(encodeURIComponent).join('/')
}
