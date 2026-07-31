/// <reference lib="webworker" />
import { clientsClaim } from 'workbox-core'
import { precacheAndRoute, getCacheKeyForURL } from 'workbox-precaching'
import { NavigationRoute, registerRoute } from 'workbox-routing'

declare const self: ServiceWorkerGlobalScope & { __WB_MANIFEST: any }

self.skipWaiting()
clientsClaim()

// Precache the built assets injected by vite-plugin-pwa.
precacheAndRoute(self.__WB_MANIFEST, {
  directoryIndex: '', // falsy
  cleanURLs: false,
})

registerRoute(
  new NavigationRoute(async ({ request }) => {
    try {
      return await fetch(request, { redirect: 'manual' })
    } catch {
      const cached = await caches.match(getCacheKeyForURL('/index.html') ?? '/index.html')
      return cached ?? Response.error()
    }
  }),
)
