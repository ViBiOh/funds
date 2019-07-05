import 'regenerator-runtime/runtime';
import { takeLatest } from 'redux-saga/effects';
import actions from 'actions';
import initSaga from './initSaga';
import getFundsSaga from './getFundsSaga';
import updateUrlSaga from './updateUrlSaga';

/**
 * Sagas of app.
 * @yield {Function} Sagas
 */
export default function* appSaga() {
  yield takeLatest(actions.INIT, initSaga);
  yield takeLatest(actions.GET_FUNDS_REQUEST, getFundsSaga);
  yield takeLatest([actions.SET_FILTER, actions.SET_ORDER, actions.SET_AGGREGAT], updateUrlSaga);
}
