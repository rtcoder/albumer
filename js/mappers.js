function getItemArtists(item) {
  return item.artists.map(artist => artist.name).join(', ');
}

function mapItem(item, content) {
  const {type, id, cover} = item;
  return `<div class="item ${type}" data-type="${type}" data-id="${id}">
           <img src="${cover || './svg/' + type + '.svg'}" alt="">
            ${content}
            <div class="close">
                <img src="./svg/close.svg" alt="">
            </div>
          </div>`;
}

function mapRelatedItem(type) {
  return (item) => {
    const {id, cover} = item;
    return `<div class="related-item ${type}" data-type="${type}" data-id="${id}">
              <img src="${cover || './svg/' + type + '.svg'}" alt="">
              <p>${item.name}</p>
          </div>`
  };
}

export function mapCd(item) {
  return mapItem(item, `
      <h2>${item.name}</h2>
      <p>${getItemArtists(item)}</p>
      <div class="related">
        ${item.artists.map(mapRelatedItem('artist')).join('')}
      </div>`
  );
}

export function mapBook(item) {
  return mapItem(item, `
      <h2>${item.name}</h2>
      <p>${getItemArtists(item)}</p>
      <div class="related">
        ${item.artists.map(mapRelatedItem('artist')).join('')}
      </div>`
  );
}

export function mapArtist(item) {
  return mapItem(item, `
      <h2>${item.name}</h2>
      <div class="related">
        ${item.albums.map(mapRelatedItem('album')).join('')}
        ${item.books.map(mapRelatedItem('book')).join('')}
      </div>
`
  );
}

