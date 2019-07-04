import funtch from 'funtch';
import { getApiUrl } from 'Constants';

/**
 * Service for dealing with Funds.
 */
export default class FundsService {
  /**
   * Retrieve funds from API.
   * @return {Array} List of funds
   */
  static async getFunds() {
    const list = await funtch.get(`${getApiUrl()}/list`);
    return list.results;
  }
}
