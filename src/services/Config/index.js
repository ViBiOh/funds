import funtch from 'funtch';
import urljoin from 'url-join';
import store from 'AppStore';

/**
 * Service for dealing with Config.
 */
export default class ConfigService {
  /**
   * Retrieve config from environment.
   * @return {Object} Config environment
   */
  static async getConfig() {
    const config = await funtch.get('/env');
    return config;
  }

  /**
   * Retrieve API URL from store
   * @param  {String} Path wanted
   * @return {String} API URL
   */
  static getApiUrl(path = '') {
    const { config: { API_URL = 'https://funds-api.vibioh.fr' } = {} } = store.getState();
    return urljoin(API_URL, path);
  }
}
