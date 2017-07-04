import { Injectable } from '@angular/core';
import { FacebookService } from 'ngx-facebook';
import {Photo, APIPhoto, FbGraphPhoto, Conflict} from '../domain/photo.model';
import { APIService } from '../api.service';
import { Channel } from '../websocket.service';
import {FriendService} from '../friends/friend.service';
import {PhotoResolver} from './photo.resolver';


@Injectable()
export class PhotoService implements Channel {

  protected photos: Map<string, Photo>;
  protected offset: string;

  constructor(
    private fb: FacebookService,
    private apiService: APIService,
    private photoResolver: PhotoResolver,
  ) {
    this.photos = new Map();
  }

  public getName(): string {
    return 'photo';
  }

  public onWebSocketMessage(data) {
    const apiPhoto = data as APIPhoto;
    const photo = this.getPhotoById(apiPhoto.id);
    this.saveToPhotoRepository(apiPhoto, photo);
  }

  public getPhotos(): Photo[] {
    return Array.from(this.photos.values());
  }

  public updatePhoto(photo: Photo) {
    this.photos.set(photo.id, photo);
    // send PUT request to API
    this.apiService.put(
      '/v1/photos/' + photo.id,
      this.photoResolver.APIPhotoFromPhoto(photo)
    ).then(response => {
      console.log(response);
    });
  }

  public getPhotoById(photoId: string): Photo {
    return this.photos.get(photoId);
  }

  public getPhotosFromFBGraph(): Promise<void> {
    let uri = '/me/photos?fields=id,created_time,from,target,images,album';

    if (this.offset) {
      uri += '&after=' + this.offset;
    }

    return this.fb.api(uri)
      .then(response => {
        if (response.paging) {
          this.offset = response.paging.cursors.after;
        }
        const fbPhotos = response.data as FbGraphPhoto[];
        return this.processFBGraphPhotos(fbPhotos);
      })
      .catch(e => console.error(e))
    ;
  }

  private processFBGraphPhotos(fbPhotos: FbGraphPhoto[]) {

    const uploaderIds = [];
    for (const fbPhoto of fbPhotos) {
      // if they have an album we can determine their privacy.
      // Add uploader id to uploaderIds to determine if they are negotiable.
      if (fbPhoto.album && !uploaderIds.includes(fbPhoto.from.id)) {
        uploaderIds.push(fbPhoto.from.id);
      }
    }

    this.apiService.get(
      '/v1/users?ids=' + JSON.stringify(uploaderIds)
    ).then(response => {
      const goodUserIds = response.json();
      const negotiablePhotos: Photo[] = [];

      for (const fbPhoto of fbPhotos) {
        if (fbPhoto.album) {
          // const photo = Photo.fromFBPhoto(fbPhoto);
          const photo = this.photoResolver.photoFromFBPhoto(fbPhoto);

          if (goodUserIds.includes(fbPhoto.from.id)) {
            photo.negotiable = true;
            negotiablePhotos.push(photo);
          }
          this.photos.set(photo.id, photo);
        }
      }
      this.updatePhotosDetail(negotiablePhotos);
    });
  }

  private updatePhotosDetail(negotiablePhotos: Photo[]) {
    const photoIds = [];
    for (const photo of negotiablePhotos) {
      photoIds.push(photo.id);
    }

    this.apiService.get(
      '/v1/photos?ids=' + JSON.stringify(photoIds)
    ).then(response => {
      const foundPhotos = response.json() as APIPhoto[];

      for (const photo of negotiablePhotos) {
        let found = false;
        for (const foundPhoto of foundPhotos) {
          if (foundPhoto.id === photo.id) {
            this.saveToPhotoRepository(foundPhoto, photo);
            found = true;
            break;
          }
        }
        if (!found) {
          // POST new photo
          this.savePhoto(photo);
        }
      }
    });
  }

  private saveToPhotoRepository(foundPhoto: APIPhoto, photo: Photo) {
    const p = this.photoResolver.photoUpdateFromAPIPhoto(photo, foundPhoto);
    this.photos.set(photo.id, p);
    console.log(p);
  }

  public savePhoto(photo: Photo) {
    this.apiService.post(
      '/v1/photos',
      photo,
    ).then(response => {
      console.log(response);
    });
  }

}
