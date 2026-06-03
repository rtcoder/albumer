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
