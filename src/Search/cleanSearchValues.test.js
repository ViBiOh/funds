import { cleanSearchValues } from './FullTextSearch';

test('should handle undefined param', () => {
  expect(cleanSearchValues()).toEqual([]);
});

test('should handle non-array @param', () => {
  expect(cleanSearchValues('test')).toEqual([]);
});

test('should not filter if under minimum value', () => {
  expect(cleanSearchValues(['test', 'unit'])).toEqual(['test', 'unit']);
});

test('should remove words below minimum', () => {
  expect(cleanSearchValues(['test', 'of', 'unit'])).toEqual(['test', 'unit']);
});

test('should not remove word if all below minimum', () => {
  expect(cleanSearchValues(['of', 'ut', 'fr'])).toEqual(['of', 'ut', 'fr']);
});
