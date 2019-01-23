import {Component, ElementRef, EventEmitter, OnInit, Output, ViewChild} from '@angular/core';
import {FormControl} from '@angular/forms';
import {SearchResultInterface} from '../../_interface/search-result.interface';
import {SearchService} from '../../_service/search.service';
import {SearchTerm} from '../../_helper/search-term';
import {debounceTime} from 'rxjs/operators';
import {MatAutocompleteSelectedEvent} from '@angular/material';
import {TypeDictionaryEnum} from '../../_enum/type-dictionary.enum';
import {Router} from '@angular/router';

@Component({
  selector: 'app-header',
  templateUrl: './header.component.html',
  styleUrls: ['./header.component.scss']
})
export class HeaderComponent implements OnInit {
  searchControl = new FormControl();
  resultsList: SearchResultInterface[];
  iconsByType = {
    genre: 'music_note',
    album: 'album',
    artist: 'person',
    group: 'people'
  };

  @ViewChild('autocompleteInput') autocompleteInput: ElementRef;
  @Output() selectedOptionChange = new EventEmitter();
  typeDictionaryEnum = TypeDictionaryEnum;

  constructor(
    private searchService: SearchService,
    private router: Router
  ) {
  }

  ngOnInit(): void {
    this.loadData();

    this.searchControl.valueChanges.pipe(
      debounceTime(300)
    ).subscribe((userInput: string) => {
      this.loadData(userInput);
    });
  }

  loadData(searchFor = ''): void {
    const search = SearchTerm.parse(searchFor);
    this.searchService.getResults(search).subscribe((results: SearchResultInterface[]) => {
      this.resultsList = results;
      console.log(this.resultsList);
    });
  }

  displayFn(post: SearchResultInterface) {
    return post ? post.name : post;
  }

  openOption($event: MatAutocompleteSelectedEvent) {
    this.searchControl.setValue('');
    this.autocompleteInput.nativeElement.blur();
    const option = $event.option.value;
    this.router.navigate([`${option.type}s/${option.id}`]);
  }
}
