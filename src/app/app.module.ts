import {BrowserModule} from '@angular/platform-browser';
import {NgModule} from '@angular/core';

import {AppRoutingModule} from './app-routing.module';
import {AppComponent} from './app.component';
import {HomeComponent} from '../components/home/home.component';
import {LeftMenuComponent} from '../components/left-menu/left-menu.component';
import {BrowserAnimationsModule} from '@angular/platform-browser/animations';
import {HeaderComponent} from '../components/header/header.component';
import {SearchService} from '../_service/search.service';
import {HTTP_INTERCEPTORS, HttpClientModule} from '@angular/common/http';
import {FormsModule, ReactiveFormsModule} from '@angular/forms';
import {MaterialModule} from './material.module';
import {PageNotFoundComponent} from '../components/page-not-found/page-not-found.component';
import {GenreListComponent} from '../components/genre/list/genre-list.component';
import {GenreNewComponent} from '../components/genre/new/genre-new.component';
import {GenreSingleComponent} from '../components/genre/single/genre-single.component';
import {LoaderComponent} from '../components/loader/loader.component';
import {LoaderInterceptor} from '../_interceptor/loader.interceptor';

@NgModule({
  declarations: [
    AppComponent,
    HomeComponent,
    LeftMenuComponent,
    HeaderComponent,
    PageNotFoundComponent,
    GenreListComponent,
    GenreNewComponent,
    GenreSingleComponent,
    LoaderComponent
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
  providers: [
    SearchService,
    {
      provide: HTTP_INTERCEPTORS,
      useClass: LoaderInterceptor,
      multi: true
    }
  ],
  bootstrap: [AppComponent]
})
export class AppModule {
}
