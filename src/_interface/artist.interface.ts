import {GroupInterface} from './group.interface';
import {AlbumInterface} from './album.interface';

export interface ArtistInterface {
  id: number;
  name: string;
  groups: GroupInterface[];
  albums: AlbumInterface[];
}
