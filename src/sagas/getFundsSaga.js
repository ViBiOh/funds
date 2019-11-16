import "regenerator-runtime/runtime";
import { call, put } from "redux-saga/effects";
import actions from "actions";
import Funds from "services/Funds";

/**
 * Saga of getFunds action
 * @yield {Function} Saga effects to sequence flow of work
 */
export default function*() {
  try {
    const funds = yield call(Funds.getFunds);
    yield put(actions.getFundsSucceeded(funds));
  } catch (e) {
    yield put(actions.getFundsFailed(e));
  }
}
