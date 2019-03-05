import {GenreInterface} from './genre.interface';
import {ArtistInterface} from './artist.interface';
import {GroupInterface} from './group.interface';

export interface AlbumInterface {
  id: number;
  name: string;
  artist: ArtistInterface | GroupInterface;
  cover: string | null;
  genre: GenreInterface;
}
