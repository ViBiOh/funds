import { replaceAccentedChar } from './index';

it('should deal with undefined param', () => {
  expect(replaceAccentedChar()).toEqual('');
});

it('should deal with null param', () => {
  expect(replaceAccentedChar(null)).toEqual('');
});

it('should remove commons french accented character', () => {
  expect(replaceAccentedChar('àéìôùÿ')).toEqual('aeiouy');
});
