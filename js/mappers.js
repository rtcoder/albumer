function getItemArtists(item) {
  return item.artists.map(artist => artist.name).join(', ');
}

function mapItem(item, content) {
  return `<div class="item ${item.type}">
           <img src="./svg/${item.type}.svg" alt="">
            ${content}
            <div class="close">
                <img src="./svg/close.svg" alt="">
            </div>
          </div>`;
}

export function mapCd(item) {
  return mapItem(item, `
      <h2>${item.name}</h2>
      <p>${getItemArtists(item)}</p>`
  );
}

export function mapBook(item) {
  return mapItem(item, `
      <h2>${item.name}</h2>
      <p>${getItemArtists(item)}</p>`
  );
}

export function mapArtist(item) {
  return mapItem(item, `
      <h2>${item.name}</h2>`
  );
}

