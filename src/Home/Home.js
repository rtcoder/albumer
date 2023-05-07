import './Home.css'
import {useEffect, useState} from "react";
import {iconsByType, sortAz} from "../helpers/helpers";
import RightPanel from "../RightPanel/RightPanel";
import {filterItems} from "../services/_db.service";
import Header from "../Header/Header";

export default function Home() {
  const [typeFilter, setTypeFilter] = useState('all');
  const [selected, setSelected] = useState(null);
  const [albums, setAlbums] = useState([]);
  const [artists, setArtists] = useState([]);
  const [books, setBooks] = useState([]);

  const isFiltered = value => {
    return ['all', value].includes(typeFilter);
  }

  const items = () => {
    return [
      ...(isFiltered('albums') ? albums : []),
      ...(isFiltered('books') ? books : []),
      ...(isFiltered('artists') ? artists : []),
    ].sort(sortAz);
  }
  const loadData = (filter = '') => {
    filterItems(filter, 'albums').then(setAlbums);
    filterItems(filter, 'books').then(setBooks);
    filterItems(filter, 'artists').then(setArtists);
  }
  const closePanel = () => {
    setSelected(null)
    items().forEach(single => single.selected = false);
  }

  const selectItem = (item) => {
    setSelected(selected?.id !== item.id ? item : null)
  }
  useEffect(() => {
    loadData()
  }, [])
  const getItemsList = () => {
    return items().map(item =>
      <div key={item.id} className="item"
           onClick={() => selectItem(item)}>
        <div className="item-content">
          <i className="material-icons">{iconsByType[item.type]}</i>
          <div className="name">{item.name}</div>
        </div>
      </div>
    )
  }
  return (
    <>
      <Header onChange={loadData}/>
      <div className="Home">
        <div className="button-toggle-group">
          <button onClick={() => setTypeFilter('all')}
                  className={`${typeFilter === 'all' ? 'selected' : ''}`}>Wszystko</button>
          <button onClick={() => setTypeFilter('albums')}
                  className={`${typeFilter === 'albums' ? 'selected' : ''}`}>Albumy</button>
          <button onClick={() => setTypeFilter('artists')}
                  className={`${typeFilter === 'artists' ? 'selected' : ''}`}>Wykonawcy</button>
          <button onClick={() => setTypeFilter('books')}
                  className={`${typeFilter === 'books' ? 'selected' : ''}`}>Książki</button>
        </div>

        <div className="grid">
          {getItemsList()}
        </div>

        <RightPanel onClosePanel={closePanel} selected={selected}/>
      </div>
    </>
  )
}
