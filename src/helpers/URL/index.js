/**
 * Retrieve search params in URL into object.
 * @return {Object} Search params
 */
// eslint-disable-next-line import/prefer-default-export
export function getSearchParamsAsObject() {
  const params = {};
  window.location.search.replace(
    /([^?&=]+)(?:=([^?&=]*))?/g,
    (match, key, value) => {
      params[key] =
        typeof value === 'undefined' ? true : decodeURIComponent(value);
    },
  );

  return params;
}
