import jsQR from 'jsqr'
import speakeasy from 'speakeasy'

export const getOtpFromImage = async (
  data: Buffer,
  width: number,
  height: number
): Promise<number> => {
  const code = jsQR(new Uint8ClampedArray(data), width, height)
  const url = new URL(code.data)
  const secret = url.searchParams.get('secret')
  const token = speakeasy.totp({
    secret: secret,
    encoding: 'base32'
  })
  return token
}
