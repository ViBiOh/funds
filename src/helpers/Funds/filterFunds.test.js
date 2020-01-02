import { filterFunds } from "./index";

it("should handle undefined funds", () => {
  expect(filterFunds()).toEqual([]);
});

it("should handle undefined filters", () => {
  expect(filterFunds([{ isin: "12345" }])).toEqual([{ isin: "12345" }]);
});

it("should filters on filters provided", () => {
  expect(
    filterFunds([{ isin: "12345" }, { isin: "67890" }], { isin: "234" })
  ).toEqual([{ isin: "12345" }]);
});

it("should handle non stringfilters", () => {
  expect(
    filterFunds([{ isin: "12345" }, { isin: "67890" }], { isin: 5 })
  ).toEqual([{ isin: "12345" }]);
});
