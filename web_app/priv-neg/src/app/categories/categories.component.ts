import {AfterViewChecked, Component, OnInit} from '@angular/core';
import { CategoryService } from './category.service';
import {Category} from '../domain/category.model';

declare var Materialize: any;

@Component({
  selector: 'app-categories',
  templateUrl: './categories.component.html'
})
export class CategoriesComponent implements OnInit, AfterViewChecked {

  newCategory: Category;
  editing: boolean;

  constructor(
    public categoryService: CategoryService,
  ) {
    this.newCategory = new Category('', true);
    this.editing = false;
  }

  ngOnInit() {
    this.categoryService.updateCategoriesFromAPI();
  }

  getGlobalCategories(): Category[] {
    const cats: Category[] = [];
    for (const cat of this.categoryService.getCategories()) {
      if (!cat.personal) {
        cats.push(cat);
      }
    }

    return cats;
  }

  getPersonalCategories(): Category[] {
    const cats: Category[] = [];
    for (const cat of this.categoryService.getCategories()) {
      if (cat.personal) {
        cats.push(cat);
      }
    }

    return cats;
  }

  toggleEdit() {
    if (!this.editing) {
      this.editing = true;
    } else {
      this.editing = false;
      if (this.newCategory.name.length > 0) {
        this.categoryService.createCategory(this.newCategory.name).catch(
          () => console.log('Invalid Category')
        ).then(
          () => this.newCategory = new Category('', true)
        );
      }
    }
  }

  ngAfterViewChecked() {
    Materialize.updateTextFields();
  }

}
