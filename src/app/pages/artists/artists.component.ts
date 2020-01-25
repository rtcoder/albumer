import {Component, OnInit} from '@angular/core';
import {ArtistService} from '../../services/artist.service';

@Component({
  selector: 'app-books',
  templateUrl: './artists.component.html',
  styleUrls: ['./artists.component.scss']
})
export class ArtistsComponent implements OnInit {

  displayedColumns: string[] = ['name'];
  dataSource: any;

  constructor(private artistService: ArtistService) {
    this.loadData();
  }

  ngOnInit() {
  }

  applyFilter(event: any) {
    this.loadData(event.target.value.trim());
  }

  loadData(search = '') {
    this.artistService.getItemsList(search).valueChanges().subscribe((data) => this.dataSource = data);
  }

}
