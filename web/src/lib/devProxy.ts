const DEFAULT_ADDR = 'localhost:8080';
const WILDCARD_HOSTS = new Set(['', '0.0.0.0', '[::]']);

export type ProxyEntry = { target: string; ws?: boolean };

// Reuses the server's listen address so both sides of the dev loop agree.
export function devProxy(env: Record<string, string | undefined>): Record<string, ProxyEntry> {
	const addr = env.POKER_HTTP_ADDR || DEFAULT_ADDR;
	const i = addr.lastIndexOf(':');
	const host = addr.slice(0, i);
	const port = addr.slice(i + 1);
	const target = `http://${WILDCARD_HOSTS.has(host) ? 'localhost' : host}:${port}`;
	return {
		'/ws': { target, ws: true },
		'/healthz': { target }
	};
}
