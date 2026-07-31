function _chunkify(t: string) {
  // Adapted from http://my.opera.com/GreyWyvern/blog/show.dml/1671288
  const tz = []
  let x = 0
  let y = -1
  let n: number | boolean = 0
  let c: string

  while (x < t.length) {
    c = t.charAt(x)
    // only include the dot in strings
    const m: boolean = (!n && c === '.') || (c >= '0' && c <= '9')
    if (m !== n) {
      // next chunk
      y++
      tz[y] = ''
      n = m
    }
    tz[y] += c
    x++
  }
  return tz
}

/**
 * Compare two strings to provide a natural sort
 * @param a first string to compare
 * @param b second string to compare
 * @return number Negative integer if b comes before a, positive integer if a comes before b
 * or 0 if the strings are identical
 */
function naturalSortCompare(a: string, b: string) {
  const aa = _chunkify(a)
  const bb = _chunkify(b)
  let x: number, aNum: any, bNum: any

  for (x = 0; aa[x] && bb[x]; x++) {
    if (aa[x] !== bb[x]) {
      aNum = Number(aa[x])
      bNum = Number(bb[x])
      // note: == is correct here

      if (aNum == aa[x] && bNum == bb[x]) {
        return aNum - bNum
      } else {
        // Forcing 'en' locale to match the server-side locale which is
        // always 'en'.
        //
        // Note: This setting isn't supported by all browsers but for the ones
        // that do there will be more consistency between client-server sorting
        return aa[x].localeCompare(bb[x], 'en')
      }
    }
  }
  return aa.length - bb.length
}

export const textUtils = {
  naturalSortCompare
}
