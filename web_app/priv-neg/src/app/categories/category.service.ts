import {Injectable} from '@angular/core';
import {APIService} from '../api.service';
import {Category} from './category.model';

@Injectable()
export class CategoryService {

  protected categories: Category[];

  constructor(
    private apiService: APIService
  ) {
    this.categories = [];
  }

  public getCategories(): Category[] {
    return this.categories;
  }

  public updateCategories(): Promise<Category[]> {
    const uri = '/v1/categories';

    return this.apiService.get(uri).then(response => {
      console.log(response);
    });
  }

  public createCategory(c: Category): Promise<Category[]> {
    const uri = '/v1/categories';

    return this.apiService.post(uri, c).then(response => {
      console.log(response);
    });
  }
}
