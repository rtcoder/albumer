import {Injectable} from '@angular/core';
import {DbService} from './db.service';
import {HttpClient} from '@angular/common/http';
import {map} from 'rxjs/operators';
import {BookInterface} from '../interfaces/book.interface';

@Injectable({
  providedIn: 'root'
})
export class BookService extends DbService {

  constructor(httpClient: HttpClient) {
    super(httpClient);
    super.getItemsList();
  }

  getItemsByFilter(filter = '') {
    return super.getItemsByFilter(filter, 'books').pipe(map((books: BookInterface[]) => {
      return books.map((book: BookInterface) => {
        book.artists = book.artistsIds.map(artistId => {
          return this._items.artists.find(artist => artist.id === artistId);
        });
        book.type = 'books';
        return book;
      });
    }));
  }
}
