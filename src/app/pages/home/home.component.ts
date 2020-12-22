import {Component, Inject} from '@angular/core';
import {AlbumService} from '../../services/album.service';
import {IconsByTypeEnum} from '../../enums/icons-by-type.enum';
import {SearchService} from '../../services/search.service';
import {AlbumInterface} from '../../interfaces/album.interface';
import {BookInterface} from '../../interfaces/book.interface';
import {ArtistInterface} from '../../interfaces/artist.interface';
import {BookService} from '../../services/book.service';
import {ArtistService} from '../../services/artist.service';
import {sortAz} from '../../helpers/helpers';
import {DataType} from '../../types/data.type';
import {DOCUMENT} from '@angular/common';

type Item = (AlbumInterface | BookInterface | ArtistInterface) & { [key: string]: any };

@Component({
  selector: 'app-home',
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.scss']
})
export class HomeComponent {
  albums: AlbumInterface[] = [];
  books: BookInterface[] = [];
  artists: ArtistInterface[] = [];

  typeFilter: 'all' | DataType = 'all';
  filter = '';

  iconsByType = IconsByTypeEnum;

  private _selected: Item;
  set selected(item: Item) {
    console.log(item);
    this._selected = item;
    if (this._selected) {
      this.document.body.classList.add('no-scroll');
    } else {
      this.document.body.classList.remove('no-scroll');
    }
  }

  get selected(): Item {
    return this._selected;
  }

  get items(): Item[] {
    return [
      ...(this.isFiltered('albums') ? this.albums : []),
      ...(this.isFiltered('books') ? this.books : []),
      ...(this.isFiltered('artists') ? this.artists : []),
    ].sort(sortAz);
  }

  constructor(private searchService: SearchService,
              private albumsService: AlbumService,
              private booksService: BookService,
              private artistsService: ArtistService,
              @Inject(DOCUMENT) private document: Document
  ) {
    searchService.search$.subscribe(value => this.loadData(value));
    this.loadData();
  }

  isFiltered(value: 'all' | DataType) {
    return ['all', value].includes(this.typeFilter);
  }

  loadData(filter = '') {
    this.filter = filter;
    this.albumsService.getItemsByFilter(filter).subscribe(items => this.albums = items);
    this.booksService.getItemsByFilter(filter).subscribe(items => this.books = items);
    this.artistsService.getItemsByFilter(filter).subscribe(items => this.artists = items);
  }

  selectItem(item: Item) {
    const newVal = !item.selected;
    this.items.forEach(single => single.selected = false);
    item.selected = newVal;
    this.selected = newVal ? item : undefined;
  }

  closePanel() {
    this.selected = undefined;
    this.items.forEach(single => single.selected = false);
  }

  search(item: Item) {
    this.searchService.search(item.name);
  }
}
