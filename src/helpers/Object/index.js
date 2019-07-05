/**
 * Safely check if given key in object is undefined
 * @param  {Object}  o   Object chcked
 * @param  {String}  key Key checked
 * @return {Boolean}     True is undefined, false otherwise
 */
export function isUndefined(o, key) {
  return !o || typeof o[key] === 'undefined';
}
