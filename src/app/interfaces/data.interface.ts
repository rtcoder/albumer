import {BookInterface} from './book.interface';
import {AlbumInterface} from './album.interface';
import {ArtistInterface} from './artist.interface';

export interface DataInterface {
  artists: ArtistInterface[];
  albums: AlbumInterface[];
  books: BookInterface[];
}
