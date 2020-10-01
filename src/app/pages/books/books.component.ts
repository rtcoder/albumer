import {Component, OnInit} from '@angular/core';
import {BookService} from '../../services/book.service';
import {ListHelper} from '../../helpers/list.helper';

@Component({
  selector: 'app-books',
  templateUrl: './books.component.html',
  styleUrls: ['./books.component.scss']
})
export class BooksComponent extends ListHelper implements OnInit {
  displayedColumns: string[] = ['cover', 'artists', 'name', 'actions'];
  service = this.bookService;

  constructor(private bookService: BookService) {
    super();
    this.loadData();
  }

  ngOnInit() {
  }
}
