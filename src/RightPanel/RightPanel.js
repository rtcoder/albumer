import './RightPanel.css'
import {iconsByType} from "../helpers/helpers";

export default function RightPanel(props) {
  return (
    <div className={`right-panel ${props.selected ? 'open' : ''}`}>
      <div className="top">
        <i className="material-icons" onClick={props.onClosePanel}>close</i>
        <div className="name">{props.selected?.name}</div>
      </div>
      <div className="middle">
        <div className="icon">
          <i className="material-icons">{iconsByType[props.selected?.type]}</i>
        </div>
        <div className="name">
          <SelectedItemArtists selected={props.selected}/>
          {props.selected?.name}
        </div>

        <SelectedItemAlbums selected={props.selected}/>
        <SelectedItemBooks selected={props.selected}/>
      </div>
    </div>
  )
}

function SelectedItemArtists(props) {
  const items = [...(props.selected?.artists || [])]
    .map(artist => artist.name).join(', ')

  return (
    <>{props.selected?.artists?.length
      ? <>{items} - </>
      : <></>
    }</>
  )
}

function SelectedItemAlbums(props) {
  const items = [...(props.selected?.albums || [])].map(album =>
    <li key={album.id}><i className="material-icons">album</i>{album.name}</li>
  )

  return (
    <>{props.selected?.albums?.length
      ? <ul className="albums">{items}</ul>
      : <></>
    }</>
  )
}

function SelectedItemBooks(props) {
  const items = [...(props.selected?.books || [])].map(book =>
    <li key={book.id}><i className="material-icons">menu_book</i>{book.name}</li>
  )

  return (
    <>{props.selected?.books?.length
      ? <ul className="books">{items}</ul>
      : <></>
    }</>
  )
}