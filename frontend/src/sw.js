/// <reference lib="webworker" />
import { clientsClaim } from 'workbox-core';
import { precacheAndRoute, getCacheKeyForURL } from 'workbox-precaching';
import { NavigationRoute, registerRoute } from 'workbox-routing';
self.skipWaiting();
clientsClaim();
// Precache the built assets injected by vite-plugin-pwa.
precacheAndRoute(self.__WB_MANIFEST);
// Navigations go to the network with redirect:'manual', so Caddy's cross-origin
// 302 to auth.domain becomes an opaqueredirect the browser follows natively
// (no more cached-shell short-circuit / reload loop). Fall back to the precached
// shell ONLY on a real network failure (genuinely offline).
registerRoute(new NavigationRoute(async ({ request }) => {
    try {
        return await fetch(request, { redirect: 'manual' });
    }
    catch {
        const cached = await caches.match(getCacheKeyForURL('/index.html') ?? '/index.html');
        return cached ?? Response.error();
    }
}));
