import { FederatedConnection, FederatedUser } from '@ownclouders/web-pkg'

export const getConnectionId = (user: FederatedUser) => {
  return `${user.user_id}@${user.idp}`
}

export const buildConnection = (user: FederatedUser): FederatedConnection => {
  return {
    id: getConnectionId(user),
    ...user
  }
}
