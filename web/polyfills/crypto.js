import { pbkdf2 } from '@noble/hashes/pbkdf2.js'
import { sha512 } from '@noble/hashes/sha2.js'

export const pbkdf2Sync = (password, salt, c, dkLen) => {
  return Buffer.from(pbkdf2(sha512, password, salt, { c, dkLen }))
}
