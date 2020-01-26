import {Component, OnInit} from '@angular/core';
import {ArtistService} from '../../services/artist.service';
import {ListHelper} from '../../helpers/list.helper';

@Component({
  selector: 'app-books',
  templateUrl: './artists.component.html',
  styleUrls: ['./artists.component.scss']
})
export class ArtistsComponent extends ListHelper implements OnInit {
  displayedColumns: string[] = ['name', 'actions'];
  service = this.artistService;

  constructor(private artistService: ArtistService) {
    super();
    this.loadData();
  }

  ngOnInit() {
  }
}
