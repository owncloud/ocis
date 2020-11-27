import encoding from 'k6/encoding';
import * as types from "../types";

export const headersDefault = ({credential}: { credential: types.Account | types.Token }): { [key: string]: string } => {
    const isOIDCGuard = (credential as types.Token).tokenType !== undefined;
    const authOIDC = credential as types.Token;
    const authBasic = credential as types.Account;

    return {
        Authorization: isOIDCGuard ? `${authOIDC.tokenType} ${authOIDC.accessToken}` : `Basic ${encoding.b64encode(`${authBasic.login}:${authBasic.password}`)}`,
    }
}