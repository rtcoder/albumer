import {Injectable} from '@angular/core';
import {AngularFireDatabase} from '@angular/fire/database';
import {DbService} from './db.service';

@Injectable({
  providedIn: 'root'
})
export class BookService extends DbService {
  basePath = '/books';

  constructor(db: AngularFireDatabase) {
    super(db);
    super.getItemsList();
  }

}
