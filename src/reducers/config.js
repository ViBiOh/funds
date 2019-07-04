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
export default function(state = initialState, action) {
  switch (action.type) {
    case actions.GET_CONFIG_SUCCEEDED:
      return {
        ...initialState,
        ...action.config,
        ready: true,
      };

    default:
      return state;
  }
}
