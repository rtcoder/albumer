import {DbServiceInterface} from '../interfaces/db-service.interface';

export class ListHelper {
  service: DbServiceInterface;
  filter = '';
  dataSource: any[] = [];
  displayedColumns: string[] = ['actions'];

  applyFilter(event: any) {
    this.filter = event.target.value.trim();
    this.loadData();
  }

  loadData() {
    this.service.getItemsList(this.filter).snapshotChanges()
      .subscribe((data) => {
        this.dataSource = data.map(value => {
          return {key: value.key, ...value.payload.exportVal()};
        });
      });
  }
}
