import './Header.css'
import {useState} from "react";

export default function Header(props) {
  const [formValue, setFormValue] = useState('');
  const handleChange = ev => {
    ev.preventDefault();
    setFormValue(ev.target.value);
    props.onChange(ev.target.value);
  }
  const clearFormValue = () => {
    setFormValue('');
    props.onChange('');
  }
  const button = formValue
    ? <i className="close-icon material-icons" onClick={clearFormValue}>close</i>
    : '';

  return (
    <header>
      <div className="form">
        <input value={formValue}
               onChange={handleChange}
               placeholder="Szukaj..."
               autoComplete="off"
               type="text"/>
        {button}
      </div>
    </header>
  )
}
