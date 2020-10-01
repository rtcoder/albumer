import {StatusEnum} from '../enums/status.enum';
import {ArtistInterface} from './artist.interface';

export interface AlbumInterface {
  id: string;
  name: string;
  artists: ArtistInterface[];
  artistsIds: string[];
  cover: string | null;
  status: StatusEnum;
  url: string;
}
