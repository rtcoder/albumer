import {Injectable} from '@angular/core';
import {DbService} from './db.service';
import {HttpClient} from '@angular/common/http';
import {Observable} from 'rxjs';
import {map} from 'rxjs/operators';
import {ArtistInterface} from '../interfaces/artist.interface';

@Injectable({
  providedIn: 'root'
})
export class ArtistService extends DbService {

  constructor(httpClient: HttpClient) {
    super(httpClient);
    super.getItemsList();
  }

  getItemsByFilter(filter = '') {
    return super.getItemsByFilter(filter, 'artists')
      .pipe(map((artists: ArtistInterface[]) =>
        artists.map(artist => this.mapArtist(artist))
      ));
  }

  getItem(id: string): Observable<any> {
    return super.getItem(id, 'artists')
      .pipe(map((artist: ArtistInterface) => this.mapArtist(artist)));
  }

  private mapArtist(artist: ArtistInterface): ArtistInterface {
    artist.albums = this._items.albums.filter(album => album.artistsIds.includes(artist.id));
    artist.books = this._items.books.filter(book => book.artistsIds.includes(artist.id));
    if (artist.artistsIds) {
      artist.artists = this._items.artists.filter(artistItem => artist.artistsIds.includes(artistItem.id));
    }
    artist.groups = this._items.artists
      .filter(artistItem => artistItem.artistsIds && artistItem.artistsIds.includes(artist.id))
      .map(group => {
        group.albums = this._items.albums.filter(album => album.artistsIds.includes(group.id));
        group.artists = undefined;
        return group;
      });
    artist.type = 'artists';
    return artist;
  }
}
