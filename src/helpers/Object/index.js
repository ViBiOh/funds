/**
 * Safely check if given key in object is undefined
 * @param  {Object}  o   Object chcked
 * @param  {String}  key Key checked
 * @return {Boolean}     True is undefined, false otherwise
 */
// eslint-disable-next-line import/prefer-default-export
export function isUndefined(o, key) {
  return !o || typeof o[key] === 'undefined';
}
