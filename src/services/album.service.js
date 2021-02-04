import DbService from './db.service';


export default class AlbumService extends DbService {

  getItemsByFilter(filter = '') {
    return super.getItemsByFilter(filter, 'albums')
      .then(albums => albums.map((album) => {
          album.artists = album.artistsIds.map(artistId =>
            this._items.artists.find(artist => artist.id === artistId)
          );
          album.type = 'albums';
          return album;
        })
      );
  }

  getItem(id) {
    return super.getItem(id, 'albums');
  }
}
