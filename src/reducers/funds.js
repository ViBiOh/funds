import actions from 'actions';
import { buildFullTextRegex, fullTextRegexFilter } from 'helpers/Search';

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

function aggregateFunds(funds, aggregat) {
  if (!aggregat.key) {
    return [];
  }

  const aggregate = {};
  const size = Math.min(funds.length, aggregat.size);
  for (let i = 0; i < size; i += 1) {
    if (typeof aggregate[funds[i][aggregat.key]] === 'undefined') {
      aggregate[funds[i][aggregat.key]] = 0;
    }
    aggregate[funds[i][aggregat.key]] += 1;
  }

  const aggregated = Object.keys(aggregate).map(label => ({
    label,
    count: aggregate[label],
  }));
  aggregated.sort((o1, o2) => o2.count - o1.count);

  return aggregated;
}

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
      const all = action.funds.filter(e => Boolean(e.id));

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
