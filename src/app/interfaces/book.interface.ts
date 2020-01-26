import {DataInterface} from './data.inerface';

export interface BookInterface extends DataInterface {
  artists: string[];
  cover: string | null;
}
