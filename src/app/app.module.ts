import {BrowserModule} from '@angular/platform-browser';
import {NgModule} from '@angular/core';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import {LoaderInterceptor} from '../_interceptor/loader.interceptor';

@NgModule({
  declarations: [
    AppComponent
    LoaderComponent
  ],
  imports: [
    BrowserModule,
    AppRoutingModule
  providers: [
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
