export class ListHelper {
  service: any;
  filter = '';
  dataSource: any[] = [];
  displayedColumns: string[] = ['actions'];

  applyFilter(event: any) {
    this.filter = event.target.value.trim();
    this.loadData();
  }

  loadData() {
    this.service.getItemsByFilter(this.filter).subscribe(data => this.dataSource = data);
  }
}
