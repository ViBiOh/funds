const IS_SECURE = process.env.API_SECURE || /^https/.test(document.location.origin);
const API_HOST = process.env.API_HOST || document.location.host.replace(/funds/i, 'funds-api');

/**
 * URL for API requests
 * @type {String}
 */
export const API = `http${IS_SECURE ? 's' : ''}://${API_HOST}/`;

/**
 * URL for source URL
 * @type {String}
 */
export const MS_URL = id =>
  `http://www.morningstar.fr/fr/funds/snapshot/snapshot.aspx?tab=0&id=${id}`;
