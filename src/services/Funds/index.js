import ConfigService from 'services/Config';

/**
 * Service for dealing with Funds.
 */
export default class FundsService {
  /**
   * Retrieve funds from API.
   * @return {Array} List of funds
   */
  static async getFunds() {
    const list = await ConfigService.getClient().get('/list');
    return list.items;
  }
}
