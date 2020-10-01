import {EventEmitter, Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {Observable, of, Subscription} from 'rxjs';
import {DataInterface} from '../interfaces/data.interface';
import {map} from 'rxjs/operators';

@Injectable({
  providedIn: 'root'
})
export class DbService {
  _items: DataInterface = null;
  get items(): Observable<DataInterface> {
    if (this._items) {
      return of(this._items);
    } else {
      return this.getItemsList();
    }
  }

  addEvent: EventEmitter<any> = new EventEmitter<any>();
  updateEvent: EventEmitter<any> = new EventEmitter<any>();
  private subscription: Subscription;

  constructor(protected httpClient: HttpClient) {
  }

  getItemsList(): Observable<DataInterface> {
    return this.httpClient.get<DataInterface>('/assets/json/data.json').pipe(map(data => {
      this.setItems(data);
      return data;
    }));
  }

  private setItems(data: DataInterface) {
    this._items = data;
  }

  getItem(id: string, type: string): Observable<any> {
    return this.items.pipe(map((data: DataInterface) => {
      return data[type].find(el => el.id === id);
    }));
  }

  getItemsByFilter(filter = '', type: string): Observable<any[]> {
    filter = filter.toLowerCase().trim();
    return this.items.pipe(map((data: DataInterface) => {
      return data[type].filter(el => {
        const name = el.name.toLowerCase();
        return name.startsWith(filter) || name.endsWith(filter) || name.includes(filter);
      });
    }));
  }

  createItem(item: any): void {
    // item.nameLower = item.name.toLowerCase();
    // this.items.push(item).then(value => this.addEvent.emit(value))
    //   .catch(error => this.handleError(error));
  }

  updateItem(key: string, value: any): void {
    // value.nameLower = value.name.toLowerCase();
    // this.items.update(key, value).then(value1 => this.updateEvent.emit(value1))
    //   .catch(error => this.handleError(error));
  }

  deleteItem(key: string): void {
    // this.items.remove(key).then(data => console.log(data))
    //   .catch(error => this.handleError(error));
  }

  deleteAll(): void {
    // this.items.remove()
    //   .catch(error => this.handleError(error));
  }

  private handleError(error) {
    console.error(error);
  }
}
