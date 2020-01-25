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
import {GroupService} from '../../services/group.service';

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
    groups: []
  };
  iconsByType = IconsByTypeEnum;

  @ViewChild('autocompleteInput', {static: false}) autocompleteInput: ElementRef;
  typeDictionaryEnum = TypeDictionarySingularEnum;

  constructor(
    private albumService: AlbumService,
    private artistService: ArtistService,
    private bookService: BookService,
    private groupService: GroupService,
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

  loadData(searchFor = ''): void {
    const search = SearchTerm.parse(searchFor);
    this.albumService.getItemsList(search).snapshotChanges().subscribe(data => {
      this.resultsList.albums = data.map(value => {
        const payload = value.payload.exportVal();
        const artists = payload.artists ? Object.values(payload.artists) : [];
        const name = artists.length ? artists.join(', ') + ' - ' + payload.name : payload.name;
        return {type: 'album', key: value.key, name, cover: payload.cover};
      });
    });
    this.artistService.getItemsList(search).snapshotChanges().subscribe(data => {
      this.resultsList.artists = data.map(value => {
        return {type: 'artist', key: value.key, name: value.payload.exportVal().name};
      });
    });
    this.groupService.getItemsList(search).snapshotChanges().subscribe(data => {
      this.resultsList.groups = data.map(value => {
        return {type: 'group', key: value.key, name: value.payload.exportVal().name};
      });
    });
    this.bookService.getItemsList(search).snapshotChanges().subscribe(data => {
      this.resultsList.books = data.map(value => {
        const payload = value.payload.exportVal();
        const artists = payload.artists ? Object.values(payload.artists) : [];
        const name = artists.length ? artists.join(', ') + ' - ' + payload.name : payload.name;
        return {type: 'book', key: value.key, name, cover: payload.cover};
      });
    });
  }

  displayFn(post: any) {
    return post ? post.name : post;
  }

  openOption($event: MatAutocompleteSelectedEvent) {
    this.searchControl.setValue('');
    this.autocompleteInput.nativeElement.blur();
    const option = $event.option.value;
    this.router.navigate([`${option.type}s/${option.key}`]);
  }
}
