import './Header.css'
import * as React from "react";

export default class Header extends React.Component {
  constructor(props) {
    super(props)
    this.state = {
      formValue: '',
      headerClass: '',
    };
    this.handleChange = this.handleChange.bind(this);
    this.clearFormValue = this.clearFormValue.bind(this);
    this.handleScroll = this.handleScroll.bind(this);
  }

  componentDidMount() {
    window.addEventListener('scroll', this.handleScroll);
  }

  componentWillUnmount() {
    window.removeEventListener('scroll', this.handleScroll);
  }

  handleScroll() {
    const headerClass = window.pageYOffset >= 200 ? 'small' : '';
    this.setState({headerClass});
  }

  handleChange(event) {
    this.setState({formValue: event.target.value});
  }

  clearFormValue() {
    this.setState({formValue: ''});
  }

  render() {
    const button = this.state.formValue
      ? <i className="close-icon material-icons" onClick={this.clearFormValue}>close</i>
      : '';

    return (
      <header className={this.state.headerClass}>
        <form>
          <input value={this.state.formValue}
                 onChange={this.handleChange}
                 placeholder="Szukaj..."
                 autoComplete="off"
                 type="text"/>
          {button}
        </form>
      </header>
    )
  }
}