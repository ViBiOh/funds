import Fetch from '../Fetch/Fetch';

export const FETCH_SIZE = 8000;

const MS_URL = id => `http://www.morningstar.fr/fr/funds/snapshot/snapshot.aspx?tab=1&id=${id}`;
const ES_PATH = `https://elasticsearch.vibioh.fr/funds/morningStarId/_search?size=${FETCH_SIZE}`;
const API_PATH = 'https://funds-api.vibioh.fr/';

export default class MorningStarService {
  static getDataUrl(id) {
    return MS_URL(id);
  }

  static getIdList() {
    return Fetch.get(ES_PATH)
      .then(data => data.hits.hits.map(doc => doc._id)); // eslint-disable-line no-underscore-dangle
  }

  static getFund(morningStarId) {
    return Fetch.get(`${API_PATH}${morningStarId}`)
      .catch(error => ({
        id: morningStarId,
        error,
      }));
  }

  static getFunds(morningStarIds) {
    return Fetch.url(`${API_PATH}list`)
      .contentJson()
      .post(morningStarIds.join(','));
  }
}
