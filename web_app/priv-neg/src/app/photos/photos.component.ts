import { AuthService } from '../auth.service';
import { Component, OnInit } from '@angular/core';
import { Photo } from '../domain/photo.model';
import { PhotoService } from './photo.service';
import {Router} from '@angular/router';


@Component({
  selector: 'app-photos',
  templateUrl: './photos.component.html',
  styles: [`
    .card-image img {
     max-height: 300px;
     max-width: 100%;
     object-fit: cover;
   }
   .blur {-webkit-filter: grayscale(100%);filter: grayscale(100%);}
  `]
})
export class PhotosComponent implements OnInit {

  public lock: boolean;

  constructor(
    private photoService: PhotoService,
    private router: Router
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
      this.photoService.getPhotosFromFBGraph()
        .then(() => this.lock = false)
      ;
    }
  }

  public selectPhoto(photo: Photo) {
    this.router.navigate(['photos', photo.id]);
  }
}
