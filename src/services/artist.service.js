import DbService from './db.service';


export default class ArtistService extends DbService {
  getItemsByFilter(filter = '') {
    return super.getItemsByFilter(filter, 'artists')
      .then(artists => artists.map(artist => this.mapArtist(artist)));
  }

  getItem(id) {
    return super.getItem(id, 'artists')
      .then(artist => this.mapArtist(artist));
  }

  mapArtist(artist) {
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
