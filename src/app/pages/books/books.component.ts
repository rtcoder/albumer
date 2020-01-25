import {Component, OnInit} from '@angular/core';
import {BookService} from '../../services/book.service';

@Component({
  selector: 'app-books',
  templateUrl: './books.component.html',
  styleUrls: ['./books.component.scss']
})
export class BooksComponent implements OnInit {

  displayedColumns: string[] = ['cover', 'artists', 'name'];
  dataSource: any;

  constructor(private bookService: BookService) {
    this.loadData();
  }

  ngOnInit() {
  }

  applyFilter(event: any) {
    this.loadData(event.target.value.trim());
  }

  loadData(search = '') {
    this.bookService.getItemsList(search).valueChanges().subscribe((data) => this.dataSource = data);
  }

}
