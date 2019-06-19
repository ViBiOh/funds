import { makeApiActionCreator } from './creator';

export default {
  ...makeApiActionCreator('getFunds', [], ['funds']),
};
