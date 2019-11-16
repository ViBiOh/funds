import { call, select } from "redux-saga/effects";
import History from "AppHistory";
import {
  ORDER_PARAM,
  ASCENDING_ORDER_PARAM,
  AGGREGAT_PARAM,
  AGGREGAT_SIZE_PARAM
} from "components/Funds/Constants";
import updateUrlSaga, { fundsSelector } from "./updateUrlSaga";

it("should read state", () => {
  const iterator = updateUrlSaga();

  expect(iterator.next().value).toEqual(select(fundsSelector));
});

it("should do nothing if nothing", () => {
  const iterator = updateUrlSaga();
  iterator.next();

  expect(
    iterator.next({
      filters: {},
      order: {},
      aggregat: {}
    }).value
  ).toEqual(call(History.push, "/"));
});

it("should update filters if provided", () => {
  const iterator = updateUrlSaga();
  iterator.next();

  expect(
    iterator.next({
      filters: {
        label: "test",
        score: "8",
        ignored: false
      },
      order: {},
      aggregat: {}
    }).value
  ).toEqual(call(History.push, "/?label=test&score=8"));
});

it("should update order if provided", () => {
  const iterator = updateUrlSaga();
  iterator.next();

  expect(
    iterator.next({
      filters: {},
      order: {
        key: "score",
        descending: true
      },
      aggregat: {}
    }).value
  ).toEqual(call(History.push, `/?${ORDER_PARAM}=score`));
});

it("should indicate ascending order if provided", () => {
  const iterator = updateUrlSaga();
  iterator.next();

  expect(
    iterator.next({
      filters: {},
      order: {
        key: "score",
        descending: false
      },
      aggregat: {}
    }).value
  ).toEqual(
    call(History.push, `/?${ORDER_PARAM}=score&${ASCENDING_ORDER_PARAM}`)
  );
});

it("should update aggregat if provided", () => {
  const iterator = updateUrlSaga();
  iterator.next();

  expect(
    iterator.next({
      filters: {},
      order: {},
      aggregat: {
        key: "score",
        size: 33
      }
    }).value
  ).toEqual(
    call(History.push, `/?${AGGREGAT_PARAM}=score&${AGGREGAT_SIZE_PARAM}=33`)
  );
});

describe("fundsSelector", () => {
  it("should handle undefined state", () => {
    expect(fundsSelector()).toEqual(undefined);
  });

  it("should return underlying funds", () => {
    expect(fundsSelector({ funds: { filters: {} } })).toEqual({ filters: {} });
  });
});
