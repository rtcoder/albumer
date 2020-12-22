import {Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {Observable, of} from 'rxjs';
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
    filter = filter.toLowerCase().trim().replace(/[^a-zA-Z ]/g, '');
    return this.items.pipe(map((data: DataInterface) => {
      return data[type].filter(el => {
        const name = el.name.toLowerCase().replace(/[^a-zA-Z ]/g, '');
        return name.startsWith(filter) || name.endsWith(filter) || name.includes(filter);
      });
    }));
  }
}
