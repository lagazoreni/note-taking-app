/// <reference lib="webworker" />
import { build, files, version } from '$service-worker';

const worker = self as unknown as ServiceWorkerGlobalScope;
const CACHE = `noted-shell-${version}`;
const SHELL = [...build, ...files];

worker.addEventListener('install', (event) => {
	event.waitUntil(
		caches
			.open(CACHE)
			.then((cache) => cache.addAll(SHELL))
			.then(() => worker.skipWaiting())
	);
});

worker.addEventListener('activate', (event) => {
	event.waitUntil(
		caches
			.keys()
			.then((keys) =>
				Promise.all(keys.filter((key) => key !== CACHE).map((key) => caches.delete(key)))
			)
			.then(() => worker.clients.claim())
	);
});

worker.addEventListener('fetch', (event) => {
	const request = event.request;
	const url = new URL(request.url);
	if (request.method !== 'GET' || url.pathname.startsWith('/api/')) return;
	event.respondWith(
		(async () => {
			const cached = await caches.match(request);
			try {
				const response = await fetch(request);
				if (response.ok && url.origin === self.location.origin) {
					const cache = await caches.open(CACHE);
					await cache.put(request, response.clone());
				}
				return response;
			} catch {
				return (
					cached ??
					(await caches.match('/offline.html')) ??
					new Response('Offline', { status: 503 })
				);
			}
		})()
	);
});
