import { describe, expect, it } from 'vitest';
import { fibonacci, prefixMatches } from './cards';

describe('prefixMatches', () => {
	it('highlights both 1 and 13 for "1"', () => {
		expect(prefixMatches(fibonacci, '1')).toEqual(['1', '13']);
	});
	it('narrows to 13 for "13"', () => {
		expect(prefixMatches(fibonacci, '13')).toEqual(['13']);
	});
	it('matches nothing for empty input', () => {
		expect(prefixMatches(fibonacci, '')).toEqual([]);
	});
});
