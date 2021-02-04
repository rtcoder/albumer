export default class DbService {
  _items = null;
  get items() {
    if (this._items) {
      return new Promise(() => {
        return this._items;
      });
    } else {
      return this.getItemsList();
    }
  }

  getItemsList() {
    return fetch('./json/data.json')
      .then(response => response.text())
      .then(text => text.length ? JSON.parse(text) : [])
      .then(data => {
        this.setItems(data);
        return data;
      });
  }

  setItems(data) {
    this._items = data;
  }

  getItem(id, type) {
    return this.items.then(data => data[type].find(el => el.id === id));
  }

  async getItemsByFilter(filter = '', type) {
    filter = filter.toLowerCase().trim().replace(/[^a-zA-Z ]/g, '');
    return this.items.then(data => data[type].filter(el => {
        const name = el.name.toLowerCase().replace(/[^a-zA-Z ]/g, '');
        return name.startsWith(filter) || name.endsWith(filter) || name.includes(filter);
      })
    );
  }
}
