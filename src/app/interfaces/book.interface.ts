import {ArtistInterface} from './artist.interface';

export interface BookInterface {
  id: string;
  artists: ArtistInterface[];
  artistsIds: string[];
  cover: string | null;
  name: string;
  url: string;
}
