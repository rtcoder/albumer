export class SearchTerm {
  static parse(term: string): string {
    return term.trim()
      .replace(/%/g, '%25')
      .replace(/\+/g, '%2B')
      .replace(/#/g, '%23')
      .replace(/\^/g, '%5E');
  }
}
