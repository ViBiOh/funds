import funtch from 'funtch';
import Constants from 'Constants';

/**
 * Service for dealing with Funds.
 */
export default class FundsService {
  /**
   * Retrieve funds from API.
   * @return {Array} List of funds
   */
  static async getFunds() {
    const list = await funtch.get(`${Constants.getApiUrl()}/list`);
    return list.results;
  }
}
