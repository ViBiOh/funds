import 'regenerator-runtime/runtime';
import { call, put } from 'redux-saga/effects';
import actions from 'actions';
import Config from 'services/Config';

/**
 * Saga of for retrieving config
 * @yield {Function} Saga effects to sequence flow of work
 */
export default function* saga() {
  try {
    const config = yield call(Config.getConfig);
    yield put(actions.getConfigSucceeded(config));
  } catch (e) {
    yield put(actions.getConfigFailed(e));
  }
}
