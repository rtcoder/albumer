import {Injectable} from '@angular/core';
import {AngularFireDatabase} from '@angular/fire/database';
import {DbService} from './db.service';

@Injectable({
  providedIn: 'root'
})
export class AlbumService extends DbService {
  basePath = '/albums';

  constructor(db: AngularFireDatabase) {
    super(db);
    super.getItemsList();
  }

}
