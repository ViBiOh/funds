import actions from 'actions';
import { orderFunds, filterFunds, aggregateFunds } from 'helpers/Funds';

/**
 * Funds reducer initial state.
 * @type {Object}
 */
export const initialState = {
  all: [],
  displayed: [],
  aggregated: [],
  filters: {},
  order: {
    key: '',
    descending: true,
  },
  aggregat: {
    key: '',
    size: 0,
  },
};

/**
 * Update list of funds by applying filter, order and aggregate
 * @param  {Array} funds     Raw List of funds
 * @param  {Object} filters  Filters applied
 * @param  {Object} order    Order of funds
 * @param  {Object} aggregat Aggregation configuration
 * @return {Object}          An object containinn displayed and aggregated list
 */
function updateList(funds, filters, order, aggregat) {
  const displayed = orderFunds(filterFunds(funds, filters), order);

  return {
    displayed,
    aggregated: aggregateFunds(displayed, aggregat),
  };
}

/**
 * Funds reducer.
 * @param {String} state  Existing Funds state
 * @param {Object} action Action dispatched
 * @return {Object} New state
 */
export default function(state = initialState, action) {
  switch (action.type) {
    case actions.GET_FUNDS_SUCCEEDED:
      const all = action.funds.filter(e => Boolean(e.isin));

      return {
        ...state,
        all,
        ...updateList(all, state.filters, state.order, state.aggregat),
      };

    case actions.SET_FILTER:
      const filters = {
        ...state.filters,
        [action.name]: action.value,
      };

      return {
        ...state,
        filters,
        ...updateList(state.all, filters, state.order, state.aggregat),
      };

    case actions.SET_ORDER:
      const order = {
        key: action.order,
        descending: action.descending,
      };

      return {
        ...state,
        order,
        ...updateList(state.all, state.filters, order, state.aggregat),
      };

    case actions.SET_AGGREGAT:
      const aggregat = {
        key: action.key,
        size: action.size,
      };

      return {
        ...state,
        aggregat,
        ...updateList(state.all, state.filters, state.order, aggregat),
      };

    default:
      return state;
  }
}
