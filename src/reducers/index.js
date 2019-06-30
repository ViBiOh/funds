import { combineReducers } from 'redux';
import pending from './pending';
import funds from './funds';

/**
 * Combined reducers of app.
 * @type {Function}
 */
export default combineReducers({
  pending,
  funds,
});
