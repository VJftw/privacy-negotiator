import {Component, OnInit} from '@angular/core';
import {ActivatedRoute, Router} from '@angular/router';
import 'rxjs/add/operator/switchMap';
import {PhotoService} from './photo.service';
import {Photo} from '../domain/photo.model';
import {CategorySelection} from '../domain/category.model';
import {CategoryService} from '../categories/category.service';


@Component({
  selector: 'app-photo',
  templateUrl: './photo-detail.component.html',
})
export class PhotoDetailComponent implements OnInit {

  public photo: Photo;
  public editing = false;
  private categorySelection: Map<string, CategorySelection> = new Map();

  constructor(
    private route: ActivatedRoute,
    private photoService: PhotoService,
    private categoryService: CategoryService,
    private router: Router
  ) {}

  ngOnInit() {
    const id = this.route.snapshot.params['id'];
    this.photo = this.photoService.getPhotoById(id);

    if (!this.photo) {
      this.router.navigate(['start']);
    }

    console.log(this.photo);
    this.updateChoices();
  }

  public getCategories(): CategorySelection[] {
    return Array.from(this.categorySelection.values());
  }

  public updateChoices() {
    for (const cat of this.categoryService.getCategories()) {
      if (this.photo.categories.includes(cat)) {
        this.categorySelection.set(cat, new CategorySelection(cat, true));
      } else {
        this.categorySelection.set(cat, new CategorySelection(cat, false));
      }
    }
  }

  toggleEdit() {
    this.editing = !this.editing;
    if (!this.editing) {
      this.photo.categories = [];
      for (const categorySelection of this.getCategories()) {
        if (categorySelection.isActive) {
          this.photo.categories.push(categorySelection.name);
        }
      }
      console.log(this.photo);
      this.photoService.updatePhoto(this.photo);
    }
  }

}


