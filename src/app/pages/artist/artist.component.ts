import {Component, OnDestroy, OnInit} from '@angular/core';
import {AlbumService} from '../../services/album.service';
import {AlbumInterface} from '../../interfaces/album.interface';
import {ArtistService} from '../../services/artist.service';
import {ActivatedRoute, NavigationEnd, Router} from '@angular/router';
import {ArtistInterface} from '../../interfaces/artist.interface';
import {BookInterface} from '../../interfaces/book.interface';
import {BookService} from '../../services/book.service';

@Component({
  selector: 'app-artist',
  templateUrl: './artist.component.html',
  styleUrls: ['./artist.component.scss']
})
export class ArtistComponent implements OnInit, OnDestroy {
  albums: AlbumInterface[] = [];
  books: BookInterface[] = [];
  artist: ArtistInterface;
  mySubscription: any;

  constructor(
    private albumService: AlbumService,
    private bookService: BookService,
    private artistService: ArtistService,
    private route: ActivatedRoute,
    private router: Router
  ) {
    this.mySubscription = this.router.events.subscribe((event) => {
      if (event instanceof NavigationEnd) {
        this.router.navigated = false;
        this.loadData();
      }
    });
  }

  ngOnInit() {
    this.loadData();
  }

  loadData() {
    const key = this.route.snapshot.paramMap.get('key');
    this.artistService.getItem(key).valueChanges().subscribe(data => {
      if (!data) {
        this.router.navigateByUrl('not-found', {skipLocationChange: true});
      } else {
        this.artist = {
          key,
          groups: [],
          name: data.name
        };
        this.albumService.getItemsList().valueChanges().subscribe(albums =>
          this.albums = albums.filter(value => value.artists && value.artists.includes(this.artist.name)).map(val => {
            return {
              name: val.name,
              artists: val.artists,
              cover: val.cover,
              key: val.key,
              status: val.status
            };
          })
        );
        this.bookService.getItemsList().valueChanges().subscribe(books =>
          this.books = books.filter(value => value.artists && value.artists.includes(this.artist.name)).map(val => {
            return {
              name: val.name,
              artists: val.artists,
              cover: val.cover,
              key: val.key
            };
          })
        );
      }
    });
  }

  ngOnDestroy() {

    if (this.mySubscription) {

      this.mySubscription.unsubscribe();

    }

  }
}
