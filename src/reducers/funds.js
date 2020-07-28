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
export default function reducer(state = initialState, action) {
  let temp;

  switch (action.type) {
    case actions.GET_FUNDS_SUCCEEDED:
      temp = action.funds.filter((e) => Boolean(e.isin));

      return {
        ...state,
        all: temp,
        ...updateList(temp, state.filters, state.order, state.aggregat),
      };

    case actions.SET_FILTER:
      temp = {
        ...state.filters,
        [action.name]: action.value,
      };

      return {
        ...state,
        filters: temp,
        ...updateList(state.all, temp, state.order, state.aggregat),
      };

    case actions.SET_ORDER:
      temp = {
        key: action.order,
        descending: action.descending,
      };

      return {
        ...state,
        order: temp,
        ...updateList(state.all, state.filters, temp, state.aggregat),
      };

    case actions.SET_AGGREGAT:
      temp = {
        key: action.key,
        size: action.size,
      };

      return {
        ...state,
        aggregat: temp,
        ...updateList(state.all, state.filters, state.order, temp),
      };

    default:
      return state;
  }
}
