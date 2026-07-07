export const getQueryParam = (paramName: string): string | null => {
  const searchParams = new URLSearchParams(window.location.search)

  return searchParams.get(paramName)
}
