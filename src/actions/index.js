import { makeActionAndTypeCreator, makeApiActionCreator } from './creator';

export default {
  ...makeActionAndTypeCreator('SET_FILTER', 'setFilter', ['name', 'value']),
  ...makeActionAndTypeCreator('SET_ORDER', 'setOrder', ['order', 'descending']),
  ...makeActionAndTypeCreator('SET_AGGREGAT', 'setAggregat', ['key', 'size']),
  ...makeApiActionCreator('getConfig', [], ['config']),
  ...makeApiActionCreator('getFunds', [], ['funds']),
};
