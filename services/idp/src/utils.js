export function withClientRequestState(obj) {
  obj.state = generateState(16);

  return obj;
}

function dec2hex (dec) {
  return dec.toString(16).padStart(2, "0")
}

// generateState :: Integer -> String
function generateState (len) {
  var arr = new Uint8Array((len || 16) / 2)
  window.crypto.getRandomValues(arr)
  return Array.from(arr, dec2hex).join('')
}
