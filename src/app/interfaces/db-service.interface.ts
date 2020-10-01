import {EventEmitter} from '@angular/core';
import {Observable} from 'rxjs';
import {DataInterface} from './data.interface';

export interface DbServiceInterface {
  _items: DataInterface;

  addEvent: EventEmitter<any>;
  updateEvent: EventEmitter<any>;

  getItemsList(query: string): Observable<any>;

  getItemsByFilter(filter: string, type?: string): Observable<any[]>;

  getItem(key: string): Observable<any>;

  createItem(item: any): void;

  updateItem(key: string, value: any): void;

  deleteItem(key: string): void;

  deleteAll(): void;
}
