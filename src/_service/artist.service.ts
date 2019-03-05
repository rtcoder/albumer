import {Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {Observable, of} from 'rxjs';
import {SearchResultInterface} from '../_interface/search-result.interface';

@Injectable({
  providedIn: 'root'
})
export class ArtistService {

  searchOption = [];
  postUrl: string = 'https://jsonplaceholder.typicode.com/posts';

  constructor(private http: HttpClient) {
  }

  getResults(searchFor: string): Observable<SearchResultInterface[]> {
    return of(<SearchResultInterface[]>[
      {
        id: 1,
        name: 'adfg',
        type: 'artist'
      },
      {
        id: 1,
        name: 'adfg',
        type: 'genre'
      },
      {
        id: 1,
        name: 'adfg',
        type: 'group'
      },
      {
        id: 1,
        name: 'adfg',
        type: 'album'
      },
      {
        id: 1,
        name: 'adfg',
        type: 'album'
      },
      {
        id: 1,
        name: 'adfg',
        type: 'album'
      },
    ]);
    return this.http.get<SearchResultInterface[]>(this.postUrl);
  }
}
