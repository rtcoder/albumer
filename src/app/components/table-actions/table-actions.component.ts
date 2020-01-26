import {Component, Input, OnInit} from '@angular/core';
import {DataType} from '../../types/data.type';
import {DbServiceInterface} from '../../interfaces/db-service.interface';
import {ArtistService} from '../../services/artist.service';
import {AlbumService} from '../../services/album.service';
import {GroupService} from '../../services/group.service';
import {BookService} from '../../services/book.service';

@Component({
  selector: 'app-table-actions',
  templateUrl: './table-actions.component.html',
  styleUrls: ['./table-actions.component.scss']
})
export class TableActionsComponent implements OnInit {
  @Input() key: string;
  _type: DataType;
  service: DbServiceInterface;

  @Input() set type(val: DataType) {
    this._type = val;
    switch (val) {
      case 'album':
        this.service = this.albumService;
        break;
      case 'artist':
        this.service = this.artistService;
        break;
      case 'group':
        this.service = this.groupService;
        break;
      case 'book':
        this.service = this.bookService;
        break;

    }
  }

  get type(): DataType {
    return this._type;
  }

  constructor(private artistService: ArtistService,
              private albumService: AlbumService,
              private bookService: BookService,
              private groupService: GroupService) {
  }

  ngOnInit() {
  }

  edit() {

  }

  delete() {
    if (confirm('Czy na pewno chcesz usunąć ten element?')) {
    // console.log(this.key)
      this.service.deleteItem(this.key);
    }

  }

}
