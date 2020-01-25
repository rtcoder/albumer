import {EventEmitter, Injectable} from '@angular/core';
import {AngularFireDatabase, AngularFireList, AngularFireObject} from '@angular/fire/database';

@Injectable({
  providedIn: 'root'
})
export class DbService {
  basePath = '';
  items: AngularFireList<any[]> = null; //  list of objects
  item: AngularFireObject<any> = null; //   single object
  addEvent: EventEmitter<any> = new EventEmitter<any>();

  constructor(private db: AngularFireDatabase) {
  }

  getItemsList(query = ''): AngularFireList<any[]> {
    this.items = this.db.list(this.basePath, ref =>
      ref.orderByChild('nameLower')
        .startAt(query.toLowerCase(), 'nameLower')
        .endAt(query.toLowerCase() + '\uf8ff', 'nameLower')
    );
    return this.items;
  }

  getItem(key: string): AngularFireObject<any> {
    const itemPath = `${this.basePath}/${key}`;
    this.item = this.db.object(itemPath);
    return this.item;
  }

  createItem(item: any): void {
    item.nameLower = item.name.toLowerCase();
    this.items.push(item).then(value => this.addEvent.emit(value))
      .catch(error => this.handleError(error));
  }

  updateItem(key: string, value: any): void {
    value.nameLower = value.name.toLowerCase();
    this.items.update(key, value)
      .catch(error => this.handleError(error));
  }

  deleteItem(key: string): void {
    this.items.remove(key)
      .catch(error => this.handleError(error));
  }

  deleteAll(): void {
    this.items.remove()
      .catch(error => this.handleError(error));
  }

  countAll() {
    return this.items.valueChanges().subscribe(data => data.length);
  }

  private handleError(error) {
    console.log(error);
  }
}
