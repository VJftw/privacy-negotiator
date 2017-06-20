import {Injectable} from '@angular/core';
import {APIService} from '../api.service';
import {Category} from './category.model';

@Injectable()
export class CategoryService {

  protected categories: string[];

  constructor(
    private apiService: APIService
  ) {
    this.categories = [];
  }

  public getCategories(): string[] {
    return this.categories;
  }

  public updateCategories(): Promise<void> {
    const uri = '/v1/categories';

    return this.apiService.get(uri).then(response => {
      for (const cat of response.json()) {
        this.categories.push(cat);
      }
    });
  }

  public createCategory(c: string): boolean {
    const uri = '/v1/categories';

    if (this.categories.includes(c)) {
      console.log('Already have' + c);
      return false;
    }

    this.categories.push(c);

    this.apiService.post(uri, new Category(c)).then(response => {
      console.log(response);
    });

    return true;
  }
}
