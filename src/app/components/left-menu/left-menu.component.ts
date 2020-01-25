import {Component} from '@angular/core';
import {MatDialog} from '@angular/material';
import {AddItemDialogComponent} from '../../dialogs/add-item-dialog/add-item-dialog.component';

@Component({
  selector: 'app-left-menu',
  templateUrl: './left-menu.component.html',
  styleUrls: ['./left-menu.component.scss']
})
export class LeftMenuComponent {

  constructor(private dialog: MatDialog) {
  }

  openAddDialog() {
    this.dialog.open(AddItemDialogComponent);
  }

}
