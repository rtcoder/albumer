import {Component, OnInit} from '@angular/core';
import {ArtistService} from '../../services/artist.service';
import {ListHelper} from '../../helpers/list.helper';
import {Router} from '@angular/router';

@Component({
  selector: 'app-books',
  templateUrl: './artists.component.html',
  styleUrls: ['./artists.component.scss']
})
export class ArtistsComponent extends ListHelper implements OnInit {
  displayedColumns: string[] = ['name', 'actions'];
  service = this.artistService;

  constructor(private artistService: ArtistService,
              private router: Router) {
    super();
    this.loadData();
  }

  ngOnInit() {
  }

  rowClick(row: any) {
    this.router.navigate([`artists/${row.key}`]);
  }
}
