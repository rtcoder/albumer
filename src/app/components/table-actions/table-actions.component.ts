import {Component, Input, OnInit} from '@angular/core';
import {DataType} from '../../types/data.type';
import {ArtistService} from '../../services/artist.service';
import {AlbumService} from '../../services/album.service';
import {BookService} from '../../services/book.service';
import {AddItemDialogComponent} from '../../dialogs/add-item-dialog/add-item-dialog.component';
import {MatDialog} from '@angular/material';

@Component({
  selector: 'app-table-actions',
  templateUrl: './table-actions.component.html',
  styleUrls: ['./table-actions.component.scss']
})
export class TableActionsComponent implements OnInit {
  @Input() key: string;
  service: any;

  constructor(private artistService: ArtistService,
              private albumService: AlbumService,
              private bookService: BookService,
              private dialog: MatDialog) {
  }

  _type: DataType;

  get type(): DataType {
    return this._type;
  }

  @Input() set type(val: DataType) {
    this._type = val;
    switch (val) {
      case 'album':
        this.service = this.albumService;
        break;
      case 'artist':
        this.service = this.artistService;
        break;
      case 'book':
        this.service = this.bookService;
        break;

    }
  }

  ngOnInit() {
  }

  edit() {
    this.dialog.open(AddItemDialogComponent, {
      width: '300px',
      data: {
        key: this.key,
        type: this.type
      }
    });
  }

  delete() {
    if (confirm('Czy na pewno chcesz usunąć ten element?')) {
      this.service.deleteItem(this.key);
    }
  }
}
