import { AuthService } from '../auth.service';
import { Component, OnInit } from '@angular/core';
import { FacebookService, InitParams, LoginResponse, LoginOptions } from 'ngx-facebook';
import { Photo } from './photo.model';
import { PhotoService } from './photo.service';


@Component({
  selector: 'app-photos',
  templateUrl: './photos.component.html',
  styles: [`.card-image img {
     max-height: 150px;
     max-width: 100%;
  }`],
  providers: [
    PhotoService
  ]
})
export class PhotosComponent implements OnInit {

  protected photos: Photo[];
  protected lock: boolean;

  constructor(
    private authService: AuthService,
    private fb: FacebookService,
    private photoService: PhotoService
  ) {
    this.photos = [];
    this.lock = false;
  }

  ngOnInit() {
    this.updateTaggedPhotos();
  }

  getTaggedPhotos() {
    return this.photoService.getPhotos()
  }

  updateTaggedPhotos() {
    if (!this.lock) {
      this.lock = true;
      this.photoService.updateTaggedPhotosForUser(this.authService.getUser())
        .then(() => this.lock = false);
    }
  }

}
