import {Injectable} from '@angular/core';
import {Subject} from 'rxjs';
import {distinctUntilChanged} from 'rxjs/operators';

@Injectable({
  providedIn: 'root'
})
export class SearchService {
  private _search = new Subject<string>();
  public search$ = this._search.asObservable().pipe(distinctUntilChanged());

  search(term: string) {
    this._search.next(term.trim());
  }
}
