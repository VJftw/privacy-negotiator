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
    this.getTaggedPhotos();
  }

  getTaggedPhotos() {
    if (!this.lock) {
      this.lock = true;
      this.photoService.getTaggedPhotosForUsera(
        this.authService.userId
      ).then(photos => this.updatePhotos(photos));
      this.lock = false;
    }
  }

  private updatePhotos(res) {
    console.log(res);
    for (const photo of res) {
      if (photo.album) {
        const p = Photo.fromFBPhoto(photo);
        // this.photoService.getPrivacyForAnAlbum(p.album.id).then(resp => console.log(resp));
        this.photos.push(Photo.fromFBPhoto(photo));
      }
    }
  }
}
