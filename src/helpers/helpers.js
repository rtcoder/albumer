export const onlyUnique = (value, index, self) => self.indexOf(value) === index;
export const sortAz = (a, b) => a.name.localeCompare(b.name);
export const iconsByType = {
  albums: 'album',
  artists: 'person',
  books: 'menu_book',
}