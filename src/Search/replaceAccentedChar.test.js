import { replaceAccentedChar } from './FullTextSearch';

test('should deal with undefined param', () => {
  expect(replaceAccentedChar()).toEqual('');
});

test('should deal with null param', () => {
  expect(replaceAccentedChar(null)).toEqual('');
});

test('should remove commons french accented character', () => {
  expect(replaceAccentedChar('àéìôùÿ')).toEqual('aeiouy');
});
