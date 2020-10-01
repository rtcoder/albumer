import {BrowserModule} from '@angular/platform-browser';
import {NgModule} from '@angular/core';

import {AppRoutingModule} from './app-routing.module';
import {AppComponent} from './app.component';
import {HomeComponent} from './pages/home/home.component';
import {LeftMenuComponent} from './components/left-menu/left-menu.component';
import {BrowserAnimationsModule} from '@angular/platform-browser/animations';
import {HeaderComponent} from './components/header/header.component';
import {HttpClientModule} from '@angular/common/http';
import {FormsModule, ReactiveFormsModule} from '@angular/forms';
import {MaterialModule} from './modules/material.module';
import {PageNotFoundComponent} from './pages/page-not-found/page-not-found.component';
import {AddItemDialogComponent} from './dialogs/add-item-dialog/add-item-dialog.component';
import {AlbumsComponent} from './pages/albums/albums.component';
import {BooksComponent} from './pages/books/books.component';
import {ArtistsComponent} from './pages/artists/artists.component';
import {NoDataInfoComponent} from './components/no-data-info/no-data-info.component';
import {TableActionsComponent} from './components/table-actions/table-actions.component';
import {ArtistComponent} from './pages/artist/artist.component';

@NgModule({
  declarations: [
    AppComponent,
    HomeComponent,
    LeftMenuComponent,
    HeaderComponent,
    PageNotFoundComponent,
    AddItemDialogComponent,
    AlbumsComponent,
    BooksComponent,
    ArtistsComponent,
    NoDataInfoComponent,
    TableActionsComponent,
    ArtistComponent
  ],
  imports: [
    BrowserModule,
    AppRoutingModule,
    BrowserAnimationsModule,
    HttpClientModule,
    FormsModule,
    ReactiveFormsModule,
    MaterialModule,
  ],
  entryComponents: [
    AddItemDialogComponent,
  ],
  providers: [],
  bootstrap: [AppComponent]
})
export class AppModule {
}
