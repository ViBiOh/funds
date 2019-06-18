import { buildFullTextRegex } from './index';

it('should have wildcard if value is empty', () => {
  expect(buildFullTextRegex(' ')).toEqual(new RegExp('[\\s\\S]*', 'gim'));
});

it('should build regex for all values', () => {
  expect(buildFullTextRegex('unit test dashboard')).toEqual(
    new RegExp(
      '[\\s\\S]*(unit|test|dashboard)[\\s\\S]*(?!\\1)(unit|test|dashboard)[\\s\\S]*(?!\\1|\\2)(unit|test|dashboard)[\\s\\S]*',
      'gim',
    ),
  );
});
