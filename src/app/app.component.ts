import {Component} from '@angular/core';

@Component({
  selector: 'app-root',
  template: `
    <div class="container">
      <header>
        <app-header></app-header>
      </header>

      <main>
        <router-outlet></router-outlet>
      </main>
    </div>
  `
})
export class AppComponent {
}
