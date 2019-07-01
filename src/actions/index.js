import { makeActionAndTypeCreator, makeApiActionCreator } from './creator';

export default {
  ...makeActionAndTypeCreator('SET_FILTER', 'setFilter', ['name', 'value']),
  ...makeApiActionCreator('getFunds', [], ['funds']),
};
