// Pobranie plików, które mają być zcache'owane
const cacheName = 'my-pwa-cache';
const filesToCache = [
  '/',
  '/index.html',
  '/manifest.json',
  // '/sample-image.jpg'

  '/css/style.css',
];

self.addEventListener('install', function (e) {
  e.waitUntil(
    caches.open(cacheName).then(function (cache) {
      return cache.addAll(filesToCache);
    })
  );
});

self.addEventListener('fetch', function (e) {
  e.respondWith(
    caches.match(e.request).then(function (response) {
      return response || fetch(e.request);
    })
  );
});
