import { AuthService } from '../auth.service';
import { Component, OnInit } from '@angular/core';
import { Photo } from './photo.model';
import { PhotoService } from './photo.service';


@Component({
  selector: 'app-photos',
  templateUrl: './photos.component.html',
  styles: [`
    .card-image img {
     max-height: 150px;
     max-width: 100%;
     object-fit: cover;
   }
   .blur {-webkit-filter: grayscale(100%);filter: grayscale(100%);}
  `],
  providers: [
    PhotoService
  ]
})
export class PhotosComponent implements OnInit {

  protected lock: boolean;

  constructor(
    private authService: AuthService,
    private photoService: PhotoService
  ) {
    this.lock = false;
  }

  ngOnInit() {
    this.updateTaggedPhotos();
  }

  getTaggedPhotos(): Photo[] {
    return this.photoService.getPhotos();
  }

  updateTaggedPhotos() {
    if (!this.lock) {
      this.lock = true;
      this.photoService.updateTaggedPhotosForUser(this.authService.getUser())
        .then(() => this.lock = false)
      ;
    }
  }

}
