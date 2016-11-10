import Fetch from '../Fetch';

const API_PATH = 'https://funds-api.vibioh.fr/';

export default class MorningStarService {
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
