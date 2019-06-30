import funtch from 'funtch';
import Constants from 'Constants';

export default class FundsService {
  static async getFunds() {
    const list = await funtch.get(`${Constants.getApiUrl()}/list`);
    return list.results;
  }
}
