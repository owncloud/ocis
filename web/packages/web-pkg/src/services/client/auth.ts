export interface AuthParameters {
  accessToken?: string
  publicLinkToken?: string
  publicLinkPassword?: string
}

export class Auth {
  accessToken?: string
  publicLinkToken?: string
  publicLinkPassword?: string

  constructor(params: AuthParameters = {}) {
    this.accessToken = params.accessToken
    this.publicLinkToken = params.publicLinkToken
    this.publicLinkPassword = params.publicLinkPassword
  }

  getHeaders(): Record<string, string> {
    return {
      ...(this.publicLinkToken && { 'public-token': this.publicLinkToken }),
      ...(this.publicLinkPassword && {
        Authorization:
          'Basic ' + Buffer.from(['public', this.publicLinkPassword].join(':')).toString('base64')
      }),
      ...(this.accessToken &&
        !this.publicLinkPassword && {
          Authorization: 'Bearer ' + this.accessToken
        })
    }
  }
}
