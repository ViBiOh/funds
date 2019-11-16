import reducer, { initialState } from "./error";

it("should have a default empty state", () => {
  expect(reducer(undefined, { type: "ON_CHANGE" })).toEqual(initialState);
});

it("should store given error on pattern action type", () => {
  expect(
    reducer(initialState, { type: "REQUEST_FAILED", error: "invalid" })
  ).toEqual("invalid");
});

it("should restore error on succeed", () => {
  expect(reducer("invalid", { type: "REQUEST_SUCCEEDED" })).toEqual(
    initialState
  );
});
