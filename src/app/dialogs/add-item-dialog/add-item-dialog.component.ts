import {Component, OnInit} from '@angular/core';
import {AlbumService} from '../../services/album.service';
import {ArtistService} from '../../services/artist.service';
import {BookService} from '../../services/book.service';
import {DataType} from '../../types/data.type';
import {TypeDictionarySingularEnum} from '../../enums/type-dictionary.enum';
import {BookInterface} from '../../interfaces/book.interface';
import {ArtistInterface} from '../../interfaces/artist.interface';
import {AlbumInterface} from '../../interfaces/album.interface';
import {GroupService} from '../../services/group.service';
import {GroupInterface} from '../../interfaces/group.interface';
import {MatAutocompleteSelectedEvent, MatDialogRef, MatSnackBar} from '@angular/material';

@Component({
  selector: 'app-add-item-dialog',
  templateUrl: './add-item-dialog.component.html',
  styleUrls: ['./add-item-dialog.component.scss']
})
export class AddItemDialogComponent implements OnInit {
  name = '';
  artist = '';
  type: DataType = 'album';
  artists: string[];

  typeDictionaryEnum = TypeDictionarySingularEnum;

  get typeDictionaryEnumKeys() {
    return Object.keys(this.typeDictionaryEnum);
  }

  constructor(
    private albumService: AlbumService,
    private artistService: ArtistService,
    private bookService: BookService,
    private groupService: GroupService,
    public dialogRef: MatDialogRef<AddItemDialogComponent>,
    private snackBar: MatSnackBar
  ) {
    console.log(this.typeDictionaryEnum);
  }

  ngOnInit() {
    this.getArtists();
  }

  sendData() {
    let service;
    let data;
    const artists = this.artist.split(',').map(val => val.trim());
    switch (this.type) {
      case 'book':
        service = this.bookService;
        data = {name: this.name, artists, cover: null} as BookInterface;
        break;
      case 'artist':
        service = this.artistService;
        data = {name: this.name} as ArtistInterface;
        break;
      case 'group':
        service = this.groupService;
        data = {name: this.name} as GroupInterface;
        break;
      case 'album':
        service = this.albumService;
        data = {name: this.name, artists, cover: null} as AlbumInterface;
        break;
    }

    if (service) {
      service.addEvent.subscribe(() => {
        service.addEvent.unsubscribe();
        this.snackBar.open('Dodano element', 'Ok', {
          duration: 2000
        });
        this.dialogRef.close();
      });
      service.createItem(data);
    }
  }

  displayFn(post: any) {
    return post ? post.name : post;
  }

  getArtists(event?: any) {
    const search = event && event.data ? event.data : '';
    this.artistService.getItemsList(search)
      .snapshotChanges().subscribe(data => {
      this.artists = data.map(value => value.payload.exportVal().name);
    });
  }

  selectArtistOption($event: MatAutocompleteSelectedEvent) {
    this.artist = $event.option.value;
  }
}
