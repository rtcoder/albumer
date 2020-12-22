import {Component, HostBinding, HostListener} from '@angular/core';
import {FormControl} from '@angular/forms';
import {SearchService} from '../../services/search.service';

@Component({
  selector: 'app-header',
  template: `
    <form>
      <input [formControl]="searchControl"
             placeholder="Szukaj..."
             autocomplete="off"
             type="text">
      <mat-icon class="close-icon"
                (click)="searchControl.setValue('')"
                *ngIf="searchControl.value">close
      </mat-icon>
    </form>
  `,
  styleUrls: ['./header.component.scss']
})
export class HeaderComponent {
  searchControl = new FormControl();

  constructor(private searchService: SearchService) {
    this.searchControl.valueChanges.subscribe(value => searchService.search(value));
    searchService.search$.subscribe(value => this.searchControl.setValue(value));
  }

  @HostBinding('class.small')
  isSmallHeader = false;

  @HostListener('window:scroll')
  onScroll() {
    this.isSmallHeader = window.pageYOffset >= 200;
  }
}
