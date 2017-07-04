import {Injectable} from '@angular/core';
import {APIService} from '../api.service';
import {APICategory, Category} from '../domain/category.model';

@Injectable()
export class CategoryService {

  protected categories: Map<string, Category>;

  constructor(
    private apiService: APIService
  ) {
    this.categories = new Map();
  }

  public getCategories(): Category[] {
    return Array.from(this.categories.values());
  }

  public getCategory(name: string) {
    return this.categories.get(name);
  }

  public updateCategoriesFromAPI(): Promise<void> {
    const uri = '/v1/categories';

    return this.apiService.get(uri).then(response => {
      for (const cat of response.json() as APICategory[]) {
        this.categories.set(cat.name, Category.fromAPICategory(cat));
      }
    });
  }

  public createCategory(c: string): Promise<any> {
    const uri = '/v1/categories';

    return this.apiService.post(uri, new Category(c)).then(response => {
      console.log(response);
      this.categories.set(c, new Category(c, true));
    });
  }
}
