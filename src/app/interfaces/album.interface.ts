import {DataInterface} from './data.inerface';

export interface AlbumInterface extends DataInterface {
  artists: string[];
  cover: string | null;
  status: 'Posiadane' | 'Zamówione';
}
