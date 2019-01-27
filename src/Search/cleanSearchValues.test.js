import { cleanSearchValues } from './FullTextSearch';

it('should handle undefined param', () => {
  expect(cleanSearchValues()).toEqual([]);
});

it('should handle non-array @param', () => {
  expect(cleanSearchValues('test')).toEqual([]);
});

it('should not filter if under minimum value', () => {
  expect(cleanSearchValues(['test', 'unit'])).toEqual(['test', 'unit']);
});

it('should remove words below minimum', () => {
  expect(cleanSearchValues(['test', 'of', 'unit'])).toEqual(['test', 'unit']);
});

it('should not remove word if all below minimum', () => {
  expect(cleanSearchValues(['of', 'ut', 'fr'])).toEqual(['of', 'ut', 'fr']);
});
