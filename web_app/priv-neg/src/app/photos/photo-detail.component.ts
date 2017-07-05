import {Component, OnInit} from '@angular/core';
import {ActivatedRoute, Router} from '@angular/router';
import 'rxjs/add/operator/switchMap';
import {PhotoService} from './photo.service';
import {Photo} from '../domain/photo.model';
import {Category, CategorySelection} from '../domain/category.model';
import {CategoryService} from '../categories/category.service';
import {User} from '../domain/user.model';


@Component({
  selector: 'app-photo',
  templateUrl: './photo-detail.component.html',
})
export class PhotoDetailComponent implements OnInit {

  public photo: Photo;
  public editing = false;
  public showConflictHelp = false;
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
      if (Category.isIn(cat, this.photo.categories)) {
        this.categorySelection.set(cat.name, new CategorySelection(cat.name, true));
      } else {
        this.categorySelection.set(cat.name, new CategorySelection(cat.name, false));
      }
    }
  }

  toggleEdit() {
    this.editing = !this.editing;
    if (!this.editing) {
      this.photo.categories = [];
      for (const categorySelection of this.getCategories()) {
        if (categorySelection.isActive) {
          this.photo.categories.push(this.categoryService.getCategory(categorySelection.name));
        }
      }
      console.log(this.photo);
      this.photoService.updatePhoto(this.photo);
    }
  }

  isInConflict(taggedUser: User): boolean {
    for (const conflict of this.photo.conflicts) {
      for (const partyUser of conflict.parties) {
        if (partyUser.id === taggedUser.id) {
          return true;
        }
      }
    }

    return false;
  }

  toggleConflictHelp() {
    this.showConflictHelp = !this.showConflictHelp;
  }


}


