(function () {
  const menus = document.querySelectorAll('[data-add-menu]');
  if (!menus.length) return;

  document.addEventListener('click', function (event) {
    menus.forEach(function (menu) {
      if (menu.open && !menu.contains(event.target)) menu.removeAttribute('open');
    });
  });

  document.addEventListener('keydown', function (event) {
    if (event.key !== 'Escape') return;
    menus.forEach(function (menu) { menu.removeAttribute('open'); });
  });
})();

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
      const note = document.createElement('span');
      note.className = 'small-note';
      note.textContent = target.dataset.emptyText || '';
      target.replaceChildren(note);
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
      .replace(/[łŁ]/g, 'l')
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

(function () {
  const filters = document.querySelectorAll('[data-list-filter]');
  if (!filters.length) return;

  const normalize = function (value) {
    return value
      .normalize('NFD')
      .replace(/[\u0300-\u036f]/g, '')
      .replace(/[łŁ]/g, 'l')
      .replace(/[^a-zA-Z0-9]+/g, '')
      .toLowerCase();
  };

  filters.forEach(function (filter) {
    const input = filter.querySelector('[data-list-filter-input]');
    const type = filter.querySelector('[data-list-filter-type]');
    const list = filter.nextElementSibling;
    if (!input || !list || !list.matches('[data-list-filter-list]')) return;

    const items = list.querySelectorAll('[data-list-filter-item]');
    const empty = list.querySelector('[data-list-filter-empty]');

    const apply = function () {
      const query = normalize(input.value);
      const selectedType = type ? type.value : '';
      let visible = 0;

      items.forEach(function (item) {
        const label = normalize(item.dataset.filterLabel || item.textContent);
        const types = (item.dataset.filterTypes || '').split(/\s+/);
        const matchesQuery = query === '' || label.includes(query);
        const matchesType = selectedType === '' || types.includes(selectedType);
        const matches = matchesQuery && matchesType;
        item.hidden = !matches;
        if (matches) visible += 1;
      });

      if (empty) empty.hidden = visible !== 0;
    };

    input.addEventListener('input', apply);
    if (type) type.addEventListener('change', apply);
  });
})();
