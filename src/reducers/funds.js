import actions from 'actions';

/**
 * Funds reducer initial state.
 * @type {Object}
 */
export const initialState = {
  funds: [],
};

/**
 * Funds reducer.
 * @param {String} state  Existing Funds state
 * @param {Object} action Action dispatched
 * @return {Object} New state
 */
export default function(state = initialState, action) {
  switch (action.type) {
    case actions.GET_FUNDS_SUCCEEDED:
      return {
        ...state,
        funds: action.funds.filter(e => Boolean(e.id)),
      };
    default:
      return state;
  }
}
