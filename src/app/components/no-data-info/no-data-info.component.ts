import {Component, Input, OnInit} from '@angular/core';

@Component({
  selector: 'app-no-data-info',
  templateUrl: './no-data-info.component.html',
  styleUrls: ['./no-data-info.component.scss']
})
export class NoDataInfoComponent implements OnInit {
  @Input() filter = '';
  @Input() dataSource = [];

  constructor() {
  }

  ngOnInit() {
  }

}
