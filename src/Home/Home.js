import './Home.css'
import * as React from "react";
import {sortAz} from "../helpers/helpers";
import RightPanel from "../RightPanel/RightPanel";
import AlbumService from "../services/album.service";
import BookService from "../services/book.service";
import ArtistService from "../services/artist.service";

export default class Home extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      typeFilter: 'all',
      albums: [],
      books: [],
      artists: [],
      selected: null,
      iconsByType: {
        albums: 'album',
        artists: 'person',
        books: 'menu_book',
      }
    };
    this.setTypeFilter = this.setTypeFilter.bind(this);
    this.getItemsList = this.getItemsList.bind(this);
    this.isFiltered = this.isFiltered.bind(this);
    this.loadData = this.loadData.bind(this);
    this.closePanel = this.closePanel.bind(this);
  }

  componentDidMount() {
    this.albumsService = new AlbumService();
    this.booksService = new BookService();
    this.artistsService = new ArtistService();
    this.loadData();
  }

  setTypeFilter(event) {
    this.setState({typeFilter: event.target.value})
  }

  getItemsList() {
    return this.items.map(item =>
      <div key={item.id} className={`item ${item.selected ? 'selected' : ''}`}
           onClick={() => this.selectItem(item)}>
        <div className="item-content">
          <i className="material-icons">{this.state.iconsByType[item.type]}</i>
          {item.name}
        </div>
      </div>
    )
  }

  get items() {
    return [
      ...(this.isFiltered('albums') ? this.state.albums : []),
      ...(this.isFiltered('books') ? this.state.books : []),
      ...(this.isFiltered('artists') ? this.state.artists : []),
    ].sort(sortAz);
  }

  set selected(item) {
    this.setState({selected: item})
    if (item) {
      document.body.classList.add('no-scroll');
    } else {
      document.body.classList.remove('no-scroll');
    }
  }

  get selected() {
    return this.state.selected;
  }

  isFiltered(value) {
    return ['all', value].includes(this.state.typeFilter);
  }

  loadData(filter = '') {
    this.filter = filter;
    this.albumsService.getItemsByFilter(filter).then(albums => this.setState({albums}));
    this.booksService.getItemsByFilter(filter).then(books => this.setState({books}));
    this.artistsService.getItemsByFilter(filter).then(artists => this.setState({artists}));
  }

  selectItem(item) {
    const newVal = !item.selected;
    this.items.forEach(single => single.selected = false);
    item.selected = newVal;
    this.selected = newVal ? item : undefined;
  }

  closePanel() {
    this.selected = undefined;
    this.items.forEach(single => single.selected = false);
  }


  render() {
    return (
      <div className="Home">
        <div className="button-toggle-group">
          <button onClick={this.setTypeFilter}
                  className={`${this.state.typeFilter === 'all' ? 'selected' : ''}`}
                  value="all">Wszystko
          </button>
          <button onClick={this.setTypeFilter}
                  className={`${this.state.typeFilter === 'albums' ? 'selected' : ''}`}
                  value="albums">Albumy
          </button>
          <button onClick={this.setTypeFilter}
                  className={`${this.state.typeFilter === 'artists' ? 'selected' : ''}`}
                  value="artists">Wykonawcy
          </button>
          <button onClick={this.setTypeFilter}
                  className={`${this.state.typeFilter === 'books' ? 'selected' : ''}`}
                  value="books">Książki
          </button>
        </div>

        <div className="grid">
          {this.getItemsList()}
        </div>

        {/*<app-no-data-info dataSource="items" filter="filter"></app-no-data-info>*/}

        <RightPanel onClosePanel={this.closePanel} selectedItem={this.selected}/>
      </div>
    )
  }
}
