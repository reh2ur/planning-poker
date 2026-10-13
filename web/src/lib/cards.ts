export const fibonacci = ['0', '1', '2', '3', '5', '8', '13', '21', '?', 'coffee'];

export function prefixMatches(deck: string[], typed: string): string[] {
	const t = typed.toLowerCase();
	if (t === '') return [];
	return deck.filter((card) => card.toLowerCase().startsWith(t));
}
