import './ProtectLayer.css'
import * as React from "react";
import { password as savedPass } from '../pass'
import Home from "../Home/Home";
import Header from "../Header/Header";

export default class ProtectLayer extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      password: '',
      shouldEnter: false,
    };

    this.putKey = this.putKey.bind(this);
    this.backspace = this.backspace.bind(this);
    this.enter = this.enter.bind(this);
  }

  putKey(key) {
    const pass = this.state.password;
    this.setState({ password: `${pass}${key}` });
  }

  backspace() {
    const pass = this.state.password;
    this.setState({ password: pass.slice(0, -1) });
  }

  enter() {
    if (savedPass === this.state.password) {
      this.setState({ shouldEnter: true });
    } else {
      document.querySelector('.input').classList.add('error');
      setTimeout(() => {
        document.querySelector('.input').classList.remove('error');
      }, 1000);
    }
  }

  get passwordBullets() {
    return this.state.password.split('').map((char, index) =>
      <b key={index}>&#8226;</b>
    );
  }

  render() {
    if (!this.state.shouldEnter) {
      return (
        <>
          <div className="ProtectLayer">
            <div className="input"> {this.passwordBullets} </div>
            <div className="keyboard">
              <div className="row">
                <div className="key" onClick={() => this.putKey('1')}>1</div>
                <div className="key" onClick={() => this.putKey('2')}>2</div>
                <div className="key" onClick={() => this.putKey('3')}>3</div>
              </div>
              <div className="row">
                <div className="key" onClick={() => this.putKey('4')}>4</div>
                <div className="key" onClick={() => this.putKey('5')}>5</div>
                <div className="key" onClick={() => this.putKey('6')}>6</div>
              </div>
              <div className="row">
                <div className="key" onClick={() => this.putKey('7')}>7</div>
                <div className="key" onClick={() => this.putKey('8')}>8</div>
                <div className="key" onClick={() => this.putKey('9')}>9</div>
              </div>
              <div className="row">
                <div className="key" onClick={() => this.putKey('0')}>0</div>
                <div className="key" onClick={this.backspace}>
                  <span className="material-icons">backspace</span>
                </div>
                <div className="key" onClick={this.enter}>
                  <span className="material-icons">login</span>
                </div>
              </div>
            </div>
          </div>
        </>
      )
    } else {
      return (
        <>
          <Header />

          <Home />
        </>
      )
    }
  }
}