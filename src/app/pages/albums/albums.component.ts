import {Component, OnInit} from '@angular/core';
import {AlbumService} from '../../services/album.service';

@Component({
  selector: 'app-albums',
  templateUrl: './albums.component.html',
  styleUrls: ['./albums.component.scss']
})
export class AlbumsComponent implements OnInit {

  displayedColumns: string[] = ['cover', 'artists', 'name'];
  dataSource: any[] = [];

  constructor(private albumService: AlbumService) {
    this.loadData();
  }

  ngOnInit() {
  }

  applyFilter(event: any) {
    this.loadData(event.target.value.trim());
  }

  loadData(search = '') {
    this.albumService.getItemsList(search).valueChanges().subscribe((data) => this.dataSource = data);
  }

}
