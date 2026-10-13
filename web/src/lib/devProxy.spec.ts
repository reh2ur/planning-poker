import { describe, expect, it } from 'vitest';
import { devProxy } from './devProxy';

describe('devProxy', () => {
	it('defaults to localhost:8080', () => {
		expect(devProxy({})).toEqual({
			'/ws': { target: 'http://localhost:8080', ws: true },
			'/healthz': { target: 'http://localhost:8080' }
		});
	});

	it('uses the port from a bare listen address', () => {
		expect(devProxy({ POKER_HTTP_ADDR: ':9191' })['/healthz']?.target).toBe(
			'http://localhost:9191'
		);
	});

	it('uses an explicit host and port', () => {
		const proxy = devProxy({ POKER_HTTP_ADDR: 'server.example.com:9000' });
		expect(proxy['/ws']).toEqual({ target: 'http://server.example.com:9000', ws: true });
	});

	it.each(['0.0.0.0:8081', '[::]:8081'])('maps wildcard listen address %s to localhost', (addr) => {
		expect(devProxy({ POKER_HTTP_ADDR: addr })['/healthz']?.target).toBe('http://localhost:8081');
	});

	it('falls back to the default for an empty value', () => {
		expect(devProxy({ POKER_HTTP_ADDR: '' })['/ws']?.target).toBe('http://localhost:8080');
	});
});
