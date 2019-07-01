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
  order: {
    key: '',
    descending: true,
  },
};

function isUndefined(o, orderKey) {
  return !o || typeof o[orderKey] === 'undefined';
}

function filterFunds(funds, filters) {
  return Object.keys(filters).reduce((previous, filter) => {
    const regex = buildFullTextRegex(String(filters[filter]));
    return previous.filter(fund => fullTextRegexFilter(fund[filter], regex));
  }, funds.slice());
}

function orderFunds(funds, { key, descending }) {
  if (!key) {
    return funds;
  }

  const compareMultiplier = descending ? -1 : 1;

  funds.sort((o1, o2) => {
    if (isUndefined(o1, key)) {
      return -1 * compareMultiplier;
    }
    if (isUndefined(o2, key)) {
      return 1 * compareMultiplier;
    }
    if (o1[key] < o2[key]) {
      return -1 * compareMultiplier;
    }
    if (o1[key] > o2[key]) {
      return 1 * compareMultiplier;
    }
    return 0;
  });

  return funds;
}

function updateList(funds, filters, order) {
  return orderFunds(filterFunds(funds, filters), order);
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
        displayed: updateList(all, state.filters, state.order),
      };

    case actions.SET_FILTER:
      const filters = {
        ...state.filters,
        [action.name]: action.value,
      };

      return {
        ...state,
        filters,
        displayed: updateList(state.all, filters, state.order),
      };

    case actions.SET_ORDER:
      const order = {
        key: action.order,
        descending: action.descending,
      };

      return {
        ...state,
        order,
        displayed: updateList(state.all, state.filters, order),
      };
    default:
      return state;
  }
}
