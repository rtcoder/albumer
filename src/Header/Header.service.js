import {EventEmitter} from "../Utils/EventEmitter";

export class HeaderService {
  static instance;

  constructor() {
    if (HeaderService.exists) {
      return HeaderService.instance;
    }
    HeaderService.instance = this;
    HeaderService.exists = true;
    return this;
  }

  _headerClass = '';

  makeHeaderSmall() {
    if (this._headerClass === '') {
      this._headerClass = 'small';
      console.log(this._headerClass)
      EventEmitter.emit('headerClass', 'small');
    }
  }

  removeHeaderSmall() {
    if (this._headerClass === 'small') {
      this._headerClass = '';
      console.log(this._headerClass)
      EventEmitter.emit('headerClass', '');
    }
  }
}