import funtch from 'funtch';
import { API } from '../Constants';

export default class FundsService {
  static getFunds() {
    return funtch.get(`${API}list`);
  }
}
