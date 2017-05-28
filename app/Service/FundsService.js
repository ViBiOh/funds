import funtch from 'funtch';

const MS_URL = id => `http://www.morningstar.fr/fr/funds/snapshot/snapshot.aspx?tab=1&id=${id}`;
const API_PATH = `https://${document.location.host.replace(/funds/i, 'funds-api')}/`;

export default class MorningStarService {
  static getDataUrl(id) {
    return MS_URL(id);
  }

  static getFunds() {
    return funtch.get(`${API_PATH}list`);
  }
}
