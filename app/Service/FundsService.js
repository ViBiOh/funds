import funtch from 'funtch';

const IS_SECURE = process.env.API_SECURE || /^https/.test(document.location.origin);
const API_HOST = process.env.API_HOST || document.location.host.replace(/funds/i, 'funds-api');

const MS_URL = id => `http://www.morningstar.fr/fr/funds/snapshot/snapshot.aspx?tab=1&id=${id}`;
const API_PATH = `http${IS_SECURE ? 's' : ''}://${API_HOST}/`;

export default class MorningStarService {
  static getDataUrl(id) {
    return MS_URL(id);
  }

  static getFunds() {
    return funtch.get(`${API_PATH}list`);
  }
}
