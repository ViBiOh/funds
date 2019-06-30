import 'regenerator-runtime/runtime';
import { call, put, takeLatest } from 'redux-saga/effects';
import actions from 'actions';
import Funds from 'services/Funds';

/**
 * Saga of GetFunds action
 * @yield {Function} Saga effects to sequence flow of work
 */
export function* getFundsSaga() {
  try {
    const funds = yield call(Funds.getFunds);
    yield put(actions.getFundsSucceeded(funds));
  } catch (e) {
    yield put(actions.getFundsFailed(e));
  }
}

/**
 * Sagas of app.
 * @yield {Function} Sagas
 */
export default function* appSaga() {
  yield takeLatest(actions.GET_FUNDS_REQUEST, getFundsSaga);
}
