import {NgModule} from '@angular/core';

import {
  MatAutocompleteModule,
  MatBadgeModule,
  MatBottomSheetModule,
  MatButtonModule,
  MatCardModule,
  MatCheckboxModule,
  MatChipsModule,
  MatDatepickerModule,
  MatDialogModule,
  MatFormFieldModule,
  MatGridListModule,
  MatIconModule,
  MatInputModule,
  MatListModule,
  MatMenuModule,
  MatNativeDateModule,
  MatPaginatorModule,
  MatProgressBarModule,
  MatProgressSpinnerModule,
  MatRadioModule,
  MatSelectModule,
  MatSlideToggleModule,
  MatSnackBarModule,
  MatSortModule,
  MatTableModule,
  MatTabsModule,
  MatTooltipModule
} from '@angular/material';

const modules = [
  MatFormFieldModule,
  MatCardModule,
  MatTableModule,
  MatButtonModule,
  MatInputModule,
  MatCheckboxModule,
  MatSelectModule,
  MatIconModule,
  MatSnackBarModule,
  MatPaginatorModule,
  MatProgressSpinnerModule,
  MatSortModule,
  MatDialogModule,
  MatTooltipModule,
  MatListModule,
  MatTabsModule,
  MatDatepickerModule,
  MatNativeDateModule,
  MatBadgeModule,
  MatGridListModule,
  MatSlideToggleModule,
  MatRadioModule,
  MatMenuModule,
  MatAutocompleteModule,
  MatProgressBarModule,
  MatBottomSheetModule,
  MatChipsModule
];

@NgModule({
  imports: [
    modules
  ],
  exports: [
    modules
  ],
  declarations: []
})
export class MaterialModule {
}


