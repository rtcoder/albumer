import {Component, ElementRef, OnInit, ViewChild} from '@angular/core';
import {FormControl} from '@angular/forms';
import {SearchTerm} from '../../helpers/search-term';
import {debounceTime} from 'rxjs/operators';
import {MatAutocompleteSelectedEvent} from '@angular/material';
import {TypeDictionarySingularEnum} from '../../enums/type-dictionary.enum';
import {Router} from '@angular/router';
import {AlbumService} from '../../services/album.service';
import {ArtistService} from '../../services/artist.service';
import {BookService} from '../../services/book.service';
import {IconsByTypeEnum} from '../../enums/icons-by-type.enum';
import {AlbumInterface} from '../../interfaces/album.interface';
import {BookInterface} from '../../interfaces/book.interface';

@Component({
  selector: 'app-header',
  templateUrl: './header.component.html',
  styleUrls: ['./header.component.scss']
})
export class HeaderComponent implements OnInit {
  searchControl = new FormControl();
  resultsList = {
    artists: [],
    albums: [],
    books: [],
  };
  iconsByType = IconsByTypeEnum;

  @ViewChild('autocompleteInput', {static: false}) autocompleteInput: ElementRef;
  typeDictionaryEnum = TypeDictionarySingularEnum;

  constructor(
    private albumService: AlbumService,
    private artistService: ArtistService,
    private bookService: BookService,
    private router: Router
  ) {
  }

  ngOnInit(): void {
    this.loadData();

    this.searchControl.valueChanges.pipe(
      debounceTime(300)
    ).subscribe((userInput: string) => {
      this.loadData(userInput);
    });
  }

  getArtists(option: AlbumInterface | BookInterface): string {
    return option.artists ? option.artists.map(artist => artist.name).join(', ') : '';
  }

  loadData(searchFor = ''): void {
    const search = SearchTerm.parse(searchFor);
    this.albumService.getItemsByFilter(search).subscribe(data => this.resultsList.albums = data.map(item => {
      item.url = 'albums/' + item.id;
      return item;
    }));
    this.artistService.getItemsByFilter(search).subscribe(data => this.resultsList.artists = data.map(item => {
      item.url = 'artists/' + item.id;
      return item;
    }));
    this.bookService.getItemsByFilter(search).subscribe(data => this.resultsList.books = data.map(item => {
      item.url = 'books/' + item.id;
      return item;
    }));
  }

  displayFn(post: any) {
    return post ? post.name : post;
  }

  openOption($event: MatAutocompleteSelectedEvent) {
    console.log($event);
    this.searchControl.setValue('');
    this.autocompleteInput.nativeElement.blur();
    const option = $event.option.value;
    this.router.navigate([option.url]);
  }
}
