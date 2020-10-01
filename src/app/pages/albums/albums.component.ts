import {Component, OnInit} from '@angular/core';
import {AlbumService} from '../../services/album.service';
import {ListHelper} from '../../helpers/list.helper';

@Component({
  selector: 'app-albums',
  templateUrl: './albums.component.html',
  styleUrls: ['./albums.component.scss']
})
export class AlbumsComponent extends ListHelper implements OnInit {
  displayedColumns: string[] = ['cover', 'name', 'artists', 'status', 'actions'];
  service = this.albumService;

  constructor(private albumService: AlbumService) {
    super();
    this.loadData();
  }

  ngOnInit() {
  }
}
