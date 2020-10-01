import {AlbumInterface} from './album.interface';
import {BookInterface} from './book.interface';

export interface ArtistInterface {
  id: string;
  name: string;
  albums: AlbumInterface[];
  books: BookInterface[];
  artists: ArtistInterface[];
  groups: ArtistInterface[];
  artistsIds: string[];
}
