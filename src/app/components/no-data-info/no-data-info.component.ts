import {Component, Input} from '@angular/core';

@Component({
  selector: 'app-no-data-info',
  template: `
    <div *ngIf="!dataSource.length" class="no-data">
      Brak danych{{filter.length ? ' spełniających podane kryteria' : ''}}
      <mat-icon>sentiment_very_dissatisfied</mat-icon>
    </div>`,
  styles: [`
    .no-data {
      display: flex;
      justify-content: center;
      color: #fefefe;
      flex-direction: column;
      align-items: center;
      padding: 30px;
      font-size: 30px;
      text-align: center;
    }

    mat-icon {
      font-size: 40px;
      width: 40px;
      height: 40px;
      margin-top: 15px;
    }
  `]
})
export class NoDataInfoComponent {
  @Input() filter = '';
  @Input() dataSource = [];
}
