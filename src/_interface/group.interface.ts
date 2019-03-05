import {ArtistInterface} from './artist.interface';
import {AlbumInterface} from './album.interface';

export interface GroupInterface {
  id: number;
  name: string;
  artists: ArtistInterface[];
  albums: AlbumInterface[];
}
