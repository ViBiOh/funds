import { isUndefined } from './index';

it('should return false is no object', () => {
  expect(isUndefined()).toEqual(true);
});

it('should return false is not an object', () => {
  expect(isUndefined(null, 'test')).toEqual(true);
  expect(isUndefined([], 'test')).toEqual(true);
  expect(isUndefined('test', 'test')).toEqual(true);
  expect(isUndefined(8000, 'test')).toEqual(true);
});

it('should return true is key is undefined', () => {
  expect(isUndefined({}, 'test')).toEqual(true);
  expect(isUndefined({ test: undefined }, 'test')).toEqual(true);
});

it('should return false is key is not undefined', () => {
  expect(isUndefined({ test: null }, 'test')).toEqual(false);
  expect(isUndefined({ test: 8000 }, 'test')).toEqual(false);
  expect(isUndefined({ test: false }, 'test')).toEqual(false);
});
