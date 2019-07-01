import actions from 'actions';
import { buildFullTextRegex, fullTextRegexFilter } from 'helpers/Search';

/**
 * Funds reducer initial state.
 * @type {Object}
 */
export const initialState = {
  all: [],
  displayed: [],
  filters: {},
};

function filterFunds(funds, filters) {
  return Object.keys(filters).reduce((previous, filter) => {
    const regex = buildFullTextRegex(String(filters[filter]));
    return previous.filter(fund => fullTextRegexFilter(fund[filter], regex));
  }, funds.slice());
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
      const all = action.funds.filter(e => Boolean(e.id));

      return {
        ...state,
        all,
        displayed: filterFunds(all, state.filters),
      };
    case actions.SET_FILTER:
      const filters = {
        ...state.filters,
        [action.name]: action.value,
      };

      return {
        ...state,
        filters,
        displayed: filterFunds(state.all, filters),
      };
    default:
      return state;
  }
}
