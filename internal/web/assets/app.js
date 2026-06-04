(function () {
  const form = document.querySelector('[data-catalog-filters]');
  const target = document.querySelector('#catalog-results');
  if (!form || !target || !window.fetch) return;

  let timer = 0;

  const load = function () {
    const params = new URLSearchParams(new FormData(form));
    const url = '/items?' + params.toString();

    window.history.replaceState(null, '', '/?' + params.toString());
    fetch(url, { headers: { 'X-Fragment': 'catalog' } })
      .then(function (res) { return res.text(); })
      .then(function (html) { target.innerHTML = html; });
  };

  form.addEventListener('input', function () {
    clearTimeout(timer);
    timer = setTimeout(load, 250);
  });
  form.addEventListener('change', load);
})();

(function () {
  const title = document.querySelector('[data-music-title]');
  const target = document.querySelector('[data-music-suggestions]');
  if (!title || !target || !window.fetch) return;

  let timer = 0;
  const load = function () {
    const query = title.value.trim();
    if (!query) {
      target.innerHTML = '<span class="small-note">Wpisz tytul, aby znalezc inne nosniki tego albumu.</span>';
      return;
    }
    const params = new URLSearchParams({ q: query, exclude: title.dataset.excludeId || '' });
    fetch('/music/suggestions?' + params.toString())
      .then(function (res) { return res.text(); })
      .then(function (html) { target.innerHTML = html; });
  };

  title.addEventListener('input', function () {
    clearTimeout(timer);
    timer = setTimeout(load, 250);
  });

  target.addEventListener('click', function (event) {
    const button = event.target.closest('[data-suggest-related]');
    if (!button) return;
    const checkbox = document.querySelector('[data-related-id="' + CSS.escape(button.dataset.suggestRelated) + '"]');
    if (checkbox) {
      checkbox.checked = true;
      checkbox.closest('.check').scrollIntoView({ block: 'nearest' });
    }
  });

  if (title.value.trim()) load();
})();

(function () {
  const groups = document.querySelectorAll('[data-check-filter]');
  if (!groups.length) return;

  const normalize = function (value) {
    return value
      .normalize('NFD')
      .replace(/[\u0300-\u036f]/g, '')
      .replace(/[^a-zA-Z0-9]+/g, '')
      .toLowerCase();
  };

  groups.forEach(function (group) {
    const input = group.querySelector('[data-check-filter-input]');
    const items = group.querySelectorAll('[data-check-filter-item]');
    if (!input || !items.length) return;

    input.addEventListener('input', function () {
      const query = normalize(input.value);
      items.forEach(function (item) {
        const label = normalize(item.dataset.filterLabel || item.textContent);
        item.hidden = query !== '' && !label.includes(query);
      });
    });
  });
})();
