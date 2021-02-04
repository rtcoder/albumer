import DbService from './db.service';

export default class BookService extends DbService {
  getItemsByFilter(filter = '') {
    return super.getItemsByFilter(filter, 'books')
      .then(books =>
        books.map((book) => {
          book.artists = book.artistsIds.map(artistId =>
            this._items.artists.find(artist => artist.id === artistId)
          );
          book.type = 'books';
          return book;
        })
      );
  }
}
