import {Component, OnInit} from '@angular/core';
import {AlbumService} from '../../services/album.service';
import {ArtistService} from '../../services/artist.service';
import {BookService} from '../../services/book.service';
import {IconsByTypeEnum} from '../../enums/icons-by-type.enum';
import {TypeDictionaryPluralEnum} from '../../enums/type-dictionary.enum';
import {StatusEnum} from '../../enums/status.enum';

@Component({
  selector: 'app-home',
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.scss']
})
export class HomeComponent implements OnInit {
  countData = {
    albums: 0,
    albums_owned: 0,
    albums_ordered: 0,
    artists: 0,
    books: 0,
  };

  typeDictionaryEnum = TypeDictionaryPluralEnum;
  iconsByType = IconsByTypeEnum;

  constructor(
    private albumService: AlbumService,
    private artistService: ArtistService,
    private bookService: BookService) {
  }

  get iconsByTypeKeys() {
    return Object.keys(this.iconsByType);
  }

  ngOnInit(): void {
    this.albumService.getItemsByFilter().subscribe(data => {
      this.countData.albums_owned = data.filter(val => val.status === StatusEnum.OWNED).length;
      this.countData.albums_ordered = data.filter(val => val.status === StatusEnum.ORDERED).length;
      this.countData.albums = data.length;
    });
    this.artistService.getItemsByFilter().subscribe(data => this.countData.artists = data.length);
    this.bookService.getItemsByFilter().subscribe(data => this.countData.books = data.length);
  }

}
