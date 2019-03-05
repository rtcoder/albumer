import {SearchResultType} from '../_types/search-result.type';

export interface SearchResultInterface {
  id: number;
  name: string;
  type: SearchResultType;
}
