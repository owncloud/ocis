/**
 * Asserts whether the given object is empty
 * @param {Object} obj Object to be checked
 * @returns {Boolean}
 */
export function isObjectEmpty (obj) {
  return Object.keys(obj).length === 0 && obj.constructor === Object
}
