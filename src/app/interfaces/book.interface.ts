import {ArtistInterface} from './artist.interface';
import {DataType} from '../types/data.type';

export interface BookInterface {
  id: string;
  artists: ArtistInterface[];
  artistsIds: string[];
  cover: string | null;
  name: string;
  url: string;
  type: DataType;
}
