import store from 'AppStore';

/**
 * Retrieve API URL from store
 * @return {String} API URL
 */
export function getApiUrl() {
  const { config: { API_URL = 'https://funds-api.vibioh.fr' } = {} } = store.getState();
  return API_URL;
}
