import { AuthService } from '../auth.service';
import { Component } from '@angular/core';
import { FacebookService, InitParams, LoginResponse, LoginOptions } from 'ngx-facebook';
import { Photo } from './photo.model';
import { PhotoService } from './photo.service';


@Component({
  selector: 'app-photos',
  templateUrl: './photos.component.html',
  providers: [
    PhotoService
  ]
})
export class PhotosComponent {

  public photos: Photo[];
  private offset: string;

  constructor(
    private authService: AuthService,
    private fb: FacebookService,
    private photoService: PhotoService
  ) {
    this.photos = [];
  }

  getTaggedPhotos() {
    if (this.offset) {
      this.photoService.getTaggedPhotosForUser(
        this.authService.userId,
        this.offset
      ).then(res => {
        this.updatePhotos(res);
      });
    } else {
      this.photoService.getTaggedPhotosForUser(
        this.authService.userId
      ).then(res => {
        this.updatePhotos(res);
      });
    }
  }

  private updatePhotos(res) {
    this.photos = this.photos.concat(res.data);
    console.log(res);
    if (res.paging) {
      this.offset = res.paging.cursors.after;
    }
  }
}
