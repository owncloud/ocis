export const extractInitials = (userName: string) => {
  return userName
    .split(/[ -]/)
    .map((part) => part.replace(/[^\p{L}\p{Nd}]/giu, ''))
    .filter(Boolean)
    .map((part) => part.charAt(0))
    .slice(0, 3)
    .join('')
    .toUpperCase()
}
