import {Component, OnDestroy, OnInit} from '@angular/core';
import {ArtistService} from '../../services/artist.service';
import {ActivatedRoute, NavigationEnd, Router} from '@angular/router';
import {ArtistInterface} from '../../interfaces/artist.interface';
import {getStatusName, StatusEnum} from '../../enums/status.enum';

@Component({
  selector: 'app-artist',
  templateUrl: './artist.component.html',
  styleUrls: ['./artist.component.scss']
})
export class ArtistComponent implements OnInit, OnDestroy {
  artists: ArtistInterface[];
  mySubscription: any;
  statusEnum = StatusEnum;
  getStatusName = getStatusName;

  constructor(
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
    const id = this.route.snapshot.paramMap.get('id');
    this.artistService.getItem(id).subscribe(data => {
      console.log(id, data);
      if (!data) {
        this.router.navigateByUrl('not-found', {skipLocationChange: true});
      } else {
        const groups = data.groups || [];
        this.artists = [data, ...groups];
      }
    });
  }

  ngOnDestroy() {
    if (this.mySubscription) {
      this.mySubscription.unsubscribe();
    }
  }
}
