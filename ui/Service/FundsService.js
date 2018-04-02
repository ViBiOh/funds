import funtch from 'funtch';
import Constants from '../Constants';

export default class FundsService {
  static getFunds() {
    return funtch.get(`${Constants.getApiUrl()}/list`);
  }
}
