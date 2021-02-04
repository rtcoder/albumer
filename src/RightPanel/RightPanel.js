import './RightPanel.css'
import * as React from "react";

export default class RightPanel extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      typeFilter: 'all',
      selected: null,
      iconsByType: {
        albums: 'album',
        artists: 'person',
        books: 'menu_book',
      }
    };
    this.closePanel = this.closePanel.bind(this);
  }

  get selected() {
    return this.props.selectedItem;
  }

  closePanel() {
    this.props.onClosePanel();
  }

  search(item) {
    // this.searchService.search(item.name);
  }

  render() {
    return (
      <div className={`right-panel ${this.selected ? 'open' : ''}`}>
        <div className="top">
          <i className="material-icons" onClick={this.closePanel}>close</i>
          <div className="name">{this.selected?.name}</div>
        </div>
        <div className="middle">
          <div className="icon">
            <i className="material-icons">{this.state.iconsByType[this.selected?.type]}</i>
          </div>
          <div className="name">
            <SelectedItemArtists selected={this.selected} onArtistClick={this.search}/>
            {this.selected?.name}
          </div>

          <SelectedItemAlbums selected={this.selected}/>
          <SelectedItemBooks selected={this.selected}/>
        </div>
      </div>
    )
  }
}

function SelectedItemArtists(props) {
  console.log(props)
  if (props.selected?.artists?.length) {
    const len = props.selected.artists.length;
    const last = (index) => index === len - 1;
    const handleClick = (artist) => {
      props.onArtistClick(artist);
    }

    const items = [...props.selected.artists].map((artist, index) =>
      <span key={artist.id} className="artist-single"
            onClick={handleClick(artist)}>
            {artist.name}{last(index) ? ' - ' : ', '}
          </span>
    )

    return (
      <>{items}</>
    )
  }
  return '';
}

function SelectedItemAlbums(props) {
  if (props.selected?.albums?.length) {
    const items = [...props.selected.albums].map(album =>
      <li key={album.id}><i className="material-icons">album</i>{album.name}</li>
    )

    return (
      <ul className="albums">{items}</ul>
    )
  }
  return '';
}

function SelectedItemBooks(props) {
  if (props.selected?.books?.length) {
    const items = [...props.selected.books].map(book =>
      <li key={book.id}><i className="material-icons">menu_book</i>{book.name}</li>
    )

    return (
      <ul className="books">{items}</ul>
    )
  }
  return '';
}