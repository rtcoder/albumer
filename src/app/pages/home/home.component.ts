import {Component, OnInit} from '@angular/core';
import {AlbumService} from '../../services/album.service';
import {ArtistService} from '../../services/artist.service';
import {BookService} from '../../services/book.service';
import {GroupService} from '../../services/group.service';
import {IconsByTypeEnum} from '../../enums/icons-by-type.enum';
import {TypeDictionaryPluralEnum} from '../../enums/type-dictionary.enum';

@Component({
  selector: 'app-home',
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.scss']
})
export class HomeComponent implements OnInit {
  countData = {
    albums: 0,
    artists: 0,
    books: 0,
    groups: 0,
  };

  typeDictionaryEnum = TypeDictionaryPluralEnum;
  iconsByType = IconsByTypeEnum;

  get iconsByTypeKeys() {
    return Object.keys(this.iconsByType);
  }

  constructor(
    private albumService: AlbumService,
    private artistService: ArtistService,
    private bookService: BookService,
    private groupService: GroupService) {
  }

  ngOnInit(): void {
    this.albumService.getItemsList().valueChanges().subscribe(data => this.countData.albums = data.length);
    this.artistService.getItemsList().valueChanges().subscribe(data => this.countData.artists = data.length);
    this.bookService.getItemsList().valueChanges().subscribe(data => this.countData.books = data.length);
    this.groupService.getItemsList().valueChanges().subscribe(data => this.countData.groups = data.length);
  }

}
