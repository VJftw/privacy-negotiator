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

  public photos: Photo[];
  private offset: string;

  constructor(
    private authService: AuthService,
    private fb: FacebookService,
    private photoService: PhotoService
  ) {
    this.photos = [];
  }

  ngOnInit() {
    this.getTaggedPhotos();
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
    console.log(res);
    // this.photos = this.photos.concat(res.data);
    for (const photo of res.data) {
      if (photo.album) {
        const p = Photo.fromGraphAPI(photo);
        this.photoService.getPrivacyForAnAlbum(p.album.id).then(resp => console.log(resp));
        this.photos.push(Photo.fromGraphAPI(photo));
      }
    }

    if (res.paging) {
      this.offset = res.paging.cursors.after;
    }
  }
}
