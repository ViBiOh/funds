import actions from 'actions';

/**
 * Config reducer initial state.
 * @type {Object}
 */
export const initialState = {
  ready: false,
};

/**
 * Config reducer.
 * @param {String} state  Existing Config state
 * @param {Object} action Action dispatched
 * @return {Object} New state
 */
export default function reducer(state = initialState, action) {
  if (action.type === actions.GET_CONFIG_SUCCEEDED) {
    return {
      ...initialState,
      ...action.config,
      ready: true,
    };
  }

  return state;
}
