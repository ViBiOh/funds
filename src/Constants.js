import funtch from 'funtch';

let context = {};

async function init() {
  const env = await funtch.get('/env');

  context = env;
  return env;
}

function getApiUrl() {
  return context.API_URL || 'https://funds-api.vibioh.fr';
}

/**
 * URL for API requests
 * @type {String}
 */
export default { init, getApiUrl };
