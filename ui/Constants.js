import funtch from 'funtch';

let context = {};

function init() {
  return new Promise((resolve) => {
    funtch.get('/env').then((env) => {
      context = env;
      resolve(context);
    });
  });
}

function getApiUrl() {
  return context.API_URL || 'https://funds-api.vibioh.fr';
}

/**
 * URL for API requests
 * @type {String}
 */
export default { init, getApiUrl };
