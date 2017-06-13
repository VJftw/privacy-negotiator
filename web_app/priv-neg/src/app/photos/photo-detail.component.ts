import {Component, OnInit} from '@angular/core';
import {ActivatedRoute, Router} from '@angular/router';
import 'rxjs/add/operator/switchMap';
import {PhotoService} from './photo.service';
import {Photo} from './photo.model';


@Component({
  selector: 'app-photo',
  templateUrl: './photo-detail.component.html',
})
export class PhotoDetailComponent implements OnInit {

  protected photo: Photo;

  autocompleteInit = {
    autocompleteOptions: {
      data: {
        'Apple': null,
        'Microsoft': null,
        'Google': null
      },
      limit: 5,
      minLength: 1
    }
  };

  constructor(
    private route: ActivatedRoute,
    private photoService: PhotoService,
    private router: Router
  ) {}

  ngOnInit() {
    const id = this.route.snapshot.params['id'];
    this.photo = this.photoService.getPhotoById(id);

    if (!this.photo) {
      this.router.navigate(['start']);
    }

    console.log(this.photo);

  }

}
