import {Injectable} from '@angular/core';
import {Subject} from 'rxjs';
import {LoaderStateInterface} from '../_interface/loader-state.interface';

@Injectable({
  providedIn: 'root'
})
export class LoaderService {

  private loaderSubject = new Subject<LoaderStateInterface>();

  loaderState = this.loaderSubject.asObservable();

  constructor() {
  }

  show() {
    this.loaderSubject.next(<LoaderStateInterface>{show: true});
  }

  hide() {
    this.loaderSubject.next(<LoaderStateInterface>{show: false});
  }

}
