export function buildWebDavPublicPath(publicLinkToken: string, path = '') {
  return `/public-files/${publicLinkToken}/${path}`.split('/').filter(Boolean).join('/')
}

export function buildWebDavOcmPath(publicLinkToken: string, path = '') {
  return `/ocm/${publicLinkToken}/${path}`.split('/').filter(Boolean).join('/')
}
