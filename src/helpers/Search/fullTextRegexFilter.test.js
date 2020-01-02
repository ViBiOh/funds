import { fullTextRegexFilter } from "./index";

it("should use given regex if provided", () => {
  expect(fullTextRegexFilter("test unit dashboard", /unit/)).toEqual(true);
});

it("should build regex if not a regex", () => {
  expect(fullTextRegexFilter("test unit dashboard", "dashboard test")).toEqual(
    true
  );
});

it("should ignore accent while matching value", () => {
  expect(fullTextRegexFilter("test unit dashboard", "dàshbôard")).toEqual(true);
});

it("should ignore accent for given value", () => {
  expect(fullTextRegexFilter("test unit dàshböard", "dashboard")).toEqual(true);
});
