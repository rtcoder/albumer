import {AngularFireList, AngularFireObject} from '@angular/fire/database';
import {EventEmitter} from '@angular/core';

export interface DbServiceInterface {
  basePath: string;
  items: AngularFireList<any[]>;
  item: AngularFireObject<any>;
  addEvent: EventEmitter<any>;
  updateEvent: EventEmitter<any>;

  getItemsList(query: string): AngularFireList<any[]>;

  getItem(key: string): AngularFireObject<any>;

  createItem(item: any): void;

  updateItem(key: string, value: any): void;

  deleteItem(key: string): void;

  deleteAll(): void;
}
