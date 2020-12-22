import {AlbumInterface} from './album.interface';
import {BookInterface} from './book.interface';
import {DataType} from '../types/data.type';

export interface ArtistInterface {
  id: string;
  name: string;
  albums: AlbumInterface[];
  books: BookInterface[];
  artists: ArtistInterface[];
  groups: ArtistInterface[];
  artistsIds: string[];
  type: DataType;
}
