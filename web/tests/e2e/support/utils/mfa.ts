import { Jimp } from 'jimp'
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

export const generateOtpFromScreenshot = async (imageBuffer: Buffer): Promise<string> => {
  const image = await Jimp.read(imageBuffer)
  const { data, width, height } = image.bitmap
  const otp = await getOtpFromImage(data, width, height)
  return otp.toString()
}
