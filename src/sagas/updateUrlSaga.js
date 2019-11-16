import "regenerator-runtime/runtime";
import { call, select } from "redux-saga/effects";
import History from "AppHistory";
import {
  ORDER_PARAM,
  ASCENDING_ORDER_PARAM,
  AGGREGAT_PARAM,
  AGGREGAT_SIZE_PARAM
} from "components/Funds/Constants";

/**
 * Selector in state for funds.
 * @param  {Object} options.funds Funds state.
 * @return {Object}               Funds state
 */
export function fundsSelector({ funds } = {}) {
  return funds;
}

/**
 * Saga of updating url from filter/agregate
 * @yield {Function} Saga effects to sequence flow of work
 */
export default function*() {
  const { filters, order, aggregat } = yield select(fundsSelector);

  const params = Object.entries(filters)
    .filter(([, value]) => Boolean(value))
    .map(([key, value]) => `${key}=${encodeURIComponent(value)}`);

  if (order.key) {
    params.push(`${ORDER_PARAM}=${encodeURIComponent(order.key)}`);

    if (!order.descending) {
      params.push(ASCENDING_ORDER_PARAM);
    }
  }

  if (aggregat.key) {
    params.push(`${AGGREGAT_PARAM}=${encodeURIComponent(aggregat.key)}`);
    params.push(`${AGGREGAT_SIZE_PARAM}=${encodeURIComponent(aggregat.size)}`);
  }

  let query = "";
  if (params.length) {
    query = `?${params.join("&")}`;
  }

  yield call(History.push, `/${query}`);
}
