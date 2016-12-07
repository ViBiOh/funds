import Fetch from '../Fetch/Fetch';

export const FETCH_SIZE = 8000;

const MS_URL = id => `http://www.morningstar.fr/fr/funds/snapshot/snapshot.aspx?tab=1&id=${id}`;
const API_PATH = 'https://funds-api.vibioh.fr/';

export default class MorningStarService {
  static getDataUrl(id) {
    return MS_URL(id);
  }

  static getFunds() {
    return Fetch.get(`${API_PATH}list`);
  }
}
