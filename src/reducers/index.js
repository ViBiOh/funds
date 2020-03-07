import { combineReducers } from 'redux';
import config from './config';
import error from './error';
import funds from './funds';
import pending from './pending';

/**
 * Combined reducers of app.
 * @type {Function}
 */
export default combineReducers({
  config,
  error,
  funds,
  pending,
});
