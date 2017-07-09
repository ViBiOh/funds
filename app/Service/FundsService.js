import funtch from 'funtch';
import { API, MS_URL } from '../Constants';

export default class MorningStarService {
  static getDataUrl(id) {
    return MS_URL(id);
  }

  static getFunds() {
    return funtch.get(`${API}list`);
  }
}
