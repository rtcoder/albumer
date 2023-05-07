const _db = {
  items: null
}
const items = () => {
  if (_db.items) {
    return new Promise((resolve) => {
      resolve(_db.items)
    });
  } else {
    return getAllItemsList();
  }
}
const setItems = items => {
  _db.items = items;
}

function getAllItemsList() {
  return fetch('./json/data.json')
    .then(response => response.text())
    .then(text => text.length ? JSON.parse(text) : [])
    .then(data => {
      setItems(data);
      return data;
    });
}

function getItemsByFilter(filter = '', type) {
  filter = filter.toLowerCase().trim().replace(/[^a-zA-Z ]/g, '');
  return items()
    .then(data => data[type].filter(el => {
        const name = el.name.toLowerCase().replace(/[^a-zA-Z ]/g, '');
        return name.startsWith(filter) || name.endsWith(filter) || name.includes(filter);
      })
    );
}

function getBooksItemsByFilter(filter = '') {
  return getItemsByFilter(filter, 'books')
    .then(books =>
      books.map((book) => {
        book.artists = book.artistsIds.map(artistId =>
          _db.items.artists.find(artist => artist.id === artistId)
        );
        book.type = 'books';
        return book;
      })
    );
}

function mapArtist(artist) {
  artist.albums = _db.items.albums.filter(album => album.artistsIds.includes(artist.id));
  artist.books = _db.items.books.filter(book => book.artistsIds.includes(artist.id));
  if (artist.artistsIds) {
    artist.artists = _db.items.artists.filter(artistItem => artist.artistsIds.includes(artistItem.id));
  }
  artist.groups = _db.items.artists
    .filter(artistItem => artistItem.artistsIds && artistItem.artistsIds.includes(artist.id))
    .map(group => {
      group.albums = _db.items.albums.filter(album => album.artistsIds.includes(group.id));
      group.artists = undefined;
      return group;
    });
  artist.type = 'artists';
  return artist;
}

function getArtistsItemsByFilter(filter = '') {
  return getItemsByFilter(filter, 'artists')
    .then(artists => artists.map(mapArtist));
}

function getAlbumsItemsByFilter(filter = '') {
  return getItemsByFilter(filter, 'albums')
    .then(albums => albums.map((album) => {
        album.artists = album.artistsIds.map(artistId =>
          _db.items.artists.find(artist => artist.id === artistId)
        );
        album.type = 'albums';
        return album;
      })
    );
}

export async function filterItems(filter = '', type = '') {
  return await {
    albums: getAlbumsItemsByFilter,
    artists: getArtistsItemsByFilter,
    books: getBooksItemsByFilter
  }[type](filter);
}
