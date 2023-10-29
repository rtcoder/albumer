import {mapArtist, mapBook, mapCd} from "./mappers.js";

document.body.addEventListener('click', e => {
  if (e.target.matches('.item')) {
    e.target.classList.add('active');
  }
  if (e.target.matches('.close')) {
    e.target.closest('.item').classList.remove('active');
  }
});


fetch('./json/data.json').then(res => res.json())
  .then(data => {
    const {books, artists, albums} = data;

    const booksList = books.map(book => {
      book.type = 'book';
      book.artists = artists.filter(artist => book.artistsIds.includes(artist.id));
      return book;
    });
    const artistsList = artists.map(artist => {
      artist.type = 'artist';
      artist.albums = albums.filter(album => album.artistsIds.includes(artist.id));
      artist.books = books.filter(book => book.artistsIds.includes(artist.id));
      return artist;
    });
    const albumsList = albums.map(album => {
      album.type = 'album';
      album.artists = artists.filter(artist => album.artistsIds.includes(artist.id));
      return album;
    });

    document.querySelector('.grid').innerHTML = [
      ...booksList,
      ...albumsList,
      ...artistsList,
    ].sort((a, b) => a.name.localeCompare(b.name))
      .map(item => {
        let val = '';
        switch (item.type) {
          case 'album':
            val = mapCd(item);
            break;
          case 'artist':
            val = mapArtist(item);
            break;
          case 'book':
            val = mapBook(item);
            break;
        }
        return val;
      }).join('');

  });