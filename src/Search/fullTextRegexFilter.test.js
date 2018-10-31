import { fullTextRegexFilter } from './FullTextSearch';

test('should use given regex if provided', () => {
  expect(fullTextRegexFilter('test unit dashboard', /unit/)).toEqual(true);
});

test('should build regex if not a regex', () => {
  expect(fullTextRegexFilter('test unit dashboard', 'dashboard test')).toEqual(true);
});

test('should ignore accent while matching value', () => {
  expect(fullTextRegexFilter('test unit dashboard', 'dàshbôard')).toEqual(true);
});

test('should ignore accent for given value', () => {
  expect(fullTextRegexFilter('test unit dàshböard', 'dashboard')).toEqual(true);
});
