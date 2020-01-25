import {Injectable} from '@angular/core';
import {AngularFireDatabase} from '@angular/fire/database';
import {DbService} from './db.service';

@Injectable({
  providedIn: 'root'
})
export class ArtistService extends DbService {
  basePath = '/artists';

  constructor(db: AngularFireDatabase) {
    super(db);
    super.getItemsList();
  }

}
