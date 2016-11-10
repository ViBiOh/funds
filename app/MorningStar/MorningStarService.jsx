import Fetch from '../Fetch';

const API_PATH = 'https://funds-api.vibioh.fr/';

export default class MorningStarService {
  static getPerformance(morningStarId) {
    return Fetch.get(`${API_PATH}${morningStarId}`)
      .catch(error => ({
        id: morningStarId,
        error,
      }));
  }

  static getPerformances(morningStarIds) {
    return Fetch.url(`${API_PATH}list`)
      .contentJson()
      .post(morningStarIds.join(','));
  }
}
