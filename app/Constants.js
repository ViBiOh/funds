const IS_SECURE = process.env.API_SECURE || /^https/.test(document.location.origin);
const API_HOST = process.env.API_HOST || document.location.host.replace(/funds/i, 'funds-api');

/**
 * URL for API requests
 * @type {String}
 */
// eslint-disable-next-line import/prefer-default-export
export const API = `http${IS_SECURE ? 's' : ''}://${API_HOST}/`;
