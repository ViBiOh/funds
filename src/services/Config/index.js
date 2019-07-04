import funtch from 'funtch';

/**
 * Service for dealing with Config.
 */
export default class ConfigService {
  /**
   * Retrieve config from environment.
   * @return {Object} Config environment
   */
  static async getFunds() {
    const config = await funtch.get('/env');
    return config;
  }
}
