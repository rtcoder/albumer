import {mapArtist, mapBook, mapCd} from "./mappers.js";

function deactivateItems() {
  document.querySelectorAll('.item').forEach(item => item.classList.remove('active'))
}

document.body.addEventListener('click', e => {
  const el = e.target;
  if (el.matches('.item')) {
    el.classList.add('active');
  }
  if (el.matches('.related-item')) {
    const {type, id} = el.dataset;
    const foundItem = document.querySelector(`.item[data-type="${type}"][data-id="${id}"]`)
    if (foundItem) {
      deactivateItems()
      foundItem.classList.add('active')
    }
  }
  if (el.matches('.close')) {
    el.closest('.item').classList.remove('active');
  }
});

fetch('http://localhost:4000/api/all').then(res => res.json())
  .then(data => {
    document.querySelector('.grid').innerHTML = data.map(el => {
      if (el.type === 'album') {
        return mapCd(el)
      }
      if (el.type === 'book') {
        return mapBook(el)
      }
      if (el.type === 'artist') {
        return mapArtist(el)
      }
    }).join('');
  });