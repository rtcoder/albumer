import {EventEmitter} from "../Utils/EventEmitter";

export class SearchService {
  _search = EventEmitter;
  search$ = this._search;

  search(term) {
    this._search.emit(term.trim());
  }
}
