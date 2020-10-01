import {Injectable} from '@angular/core';
import {DbService} from './db.service';
import {HttpClient} from '@angular/common/http';
import {map} from 'rxjs/operators';
import {AlbumInterface} from '../interfaces/album.interface';

@Injectable({
  providedIn: 'root'
})
export class AlbumService extends DbService {

  constructor(httpClient: HttpClient) {
    super(httpClient);
  }

  getItemsByFilter(filter = '') {
    return super.getItemsByFilter(filter, 'albums').pipe(map((albums: AlbumInterface[]) => {
      return albums.map((album: AlbumInterface) => {
        album.artists = album.artistsIds.map(artistId => {
          return this._items.artists.find(artist => artist.id === artistId);
        });
        return album;
      });
    }));
  }
}
