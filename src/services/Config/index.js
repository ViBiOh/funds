import funtch from 'funtch';

/**
 * Default http client.
 * @type {Object}
 */
let httpClient = funtch.withDefault({
  baseURL: 'https://funds-api.vibioh.fr',
});

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

    if (config && config.API_URL) {
      httpClient = funtch.withDefault({
        baseURL: config.API_URL,
      });
    }

    return config;
  }

  /**
   * Retrieves client with base URL for requesting API
   * @return {Object} Funtch object
   */
  static getClient() {
    return httpClient;
  }
}
