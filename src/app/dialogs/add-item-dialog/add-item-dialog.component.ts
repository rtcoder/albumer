import {Component, Inject, OnDestroy, OnInit} from '@angular/core';
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
import {MAT_DIALOG_DATA, MatAutocompleteSelectedEvent, MatDialogRef, MatSnackBar} from '@angular/material';
import {onlyUnique} from '../../helpers/helpers';
import {DbServiceInterface} from '../../interfaces/db-service.interface';

interface DialogData {
  name: string;
  type: DataType;
  key: string;
}

@Component({
  selector: 'app-add-item-dialog',
  templateUrl: './add-item-dialog.component.html',
  styleUrls: ['./add-item-dialog.component.scss']
})
export class AddItemDialogComponent implements OnInit, OnDestroy {

  subscribed = false;
  subscription;
  name = '';
  artist = '';
  artists = [];
  status = 'Posiadane';
  type: DataType = 'album';
  artistsList: string[];

  typeDictionaryEnum = TypeDictionarySingularEnum;

  get service(): DbServiceInterface {
    let service;
    switch (this.type) {
      case 'book':
        service = this.bookService;
        break;
      case 'artist':
        service = this.artistService;
        break;
      case 'group':
        service = this.groupService;
        break;
      case 'album':
        service = this.albumService;
        break;
    }
    return service;
  }

  constructor(
    private albumService: AlbumService,
    private artistService: ArtistService,
    private bookService: BookService,
    private groupService: GroupService,
    public dialogRef: MatDialogRef<AddItemDialogComponent>,
    private snackBar: MatSnackBar,
    @Inject(MAT_DIALOG_DATA) public data: DialogData
  ) {
  }

  get typeDictionaryEnumKeys() {
    return Object.keys(this.typeDictionaryEnum);
  }

  get filteredArtists(): string[] {
    return this.artists.map(val => val.trim()).filter(val => val.length).filter(onlyUnique);
  }

  ngOnInit() {
    this.getArtists();
    if (this.data.key) {
      this.type = this.data.type;
      this.service.getItem(this.data.key).valueChanges().subscribe(value => {
        console.log(value);
        this.artists = value.artists;
        this.name = value.name;
        this.status = value.status;
        this.name = value.name;
      });
    }
  }

  sendData() {
    let data;
    if (this.artist.length) {
      this.artists.push(this.artist);
    }
    const artists = this.filteredArtists;
    switch (this.type) {
      case 'book':
        data = {name: this.name, artists, cover: null} as BookInterface;
        break;
      case 'artist':
        data = {name: this.name} as ArtistInterface;
        break;
      case 'group':
        data = {name: this.name} as GroupInterface;
        break;
      case 'album':
        data = {name: this.name, artists, cover: null, status: this.status} as AlbumInterface;
        break;
    }
    if (['book', 'album'].includes(this.type)) {
      const elementsToAdd = [];
      this.artistService.getItemsList().valueChanges().subscribe(val => {
        artists.forEach(artist => {
          const found = val.find(obj => obj.name && obj.name === artist);
          if (!found) {
            elementsToAdd.push(artist);
          }
        });
      });
      setTimeout(() => elementsToAdd.forEach(val => this.artistService.createItem({name: val})), 10);
    }

    if (this.service) {
      this.subscribed = true;
      if (this.data && this.data.key) {
        console.log('here', data, this.data.key);
        this.subscription = this.service.updateEvent.subscribe(() => {
          this.snackBar.open('Zmieniono element', 'Ok', {
            duration: 2000
          });
          this.dialogRef.close();
        });
        this.service.updateItem(this.data.key, data);
      } else {
        this.subscription = this.service.addEvent.subscribe(() => {
          this.snackBar.open('Dodano element', 'Ok', {
            duration: 2000
          });
          this.dialogRef.close();
        });
        this.service.createItem(data);
      }
    }
  }

  displayFn(post: any) {
    return post ? post.name : post;
  }

  getArtists(event?: any) {
    const search = event && event.target.value ? event.target.value.trim() : '';
    this.artistService.getItemsList(search)
      .snapshotChanges().subscribe(data => {
      this.artistsList = data.map(value => value.payload.exportVal().name);
    });
  }

  selectArtistOption($event: MatAutocompleteSelectedEvent) {
    this.artists.push($event.option.value);
  }

  addArtist($event: any) {
    $event.preventDefault();
    if ($event.target.value) {
      this.artists.push($event.target.value);
      $event.target.value = '';
      this.artist = '';
    }
  }

  ngOnDestroy(): void {
    if (this.service && this.subscribed) {
      this.subscription.unsubscribe();
    }
  }
}