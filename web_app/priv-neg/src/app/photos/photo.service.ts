import { Injectable } from '@angular/core';
import { FacebookService } from 'ngx-facebook';
import {Photo, APIPhoto, FbGraphPhoto, Conflict} from '../domain/photo.model';
import { APIService } from '../api.service';
import { Channel } from '../websocket.service';
import {FriendService} from '../friends/friend.service';
import {PhotoResolver} from './photo.resolver';

class PromisePhoto {
  promise: Promise<any>;
  photo: Photo;
}

@Injectable()
export class PhotoService implements Channel {

  protected photos: Map<string, PromisePhoto>;
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
    this.getPhotoById(apiPhoto.id).then(photo => {
      const p = this.photoResolver.photoUpdateFromAPIPhoto(photo, apiPhoto);
      p.negotiable = true;
      console.log(p);
      const pP = this.photos.get(p.id);
      pP.photo = p;
      this.photos.set(p.id, pP);
    });
  }

  public getPhotos(): Photo[] {
    const photos = [];
    for (const photoPromise of Array.from(this.photos.values())) {
      if (photoPromise.photo) {
        photos.push(photoPromise.photo);
      }
    }
    return photos;
  }

  public updatePhoto(photo: Photo) {
    const p = this.photos.get(photo.id);
    p.photo = photo;
    this.photos.set(photo.id, p);
    // send PUT request to API
    this.apiService.put(
      '/v1/photos/' + photo.id,
      this.photoResolver.APIPhotoFromPhoto(photo)
    ).then(response => {
      console.log(response);
    });
  }

  public getPhotoById(id: string): Promise<Photo> {
    return new Promise((resolve, reject) => {
      if (!this.photos.has(id) || !this.photos.get(id).photo.negotiable) {
        const pP = new PromisePhoto();
        pP.promise = this.fb.api('/' + id + '?fields=id,created_time,from,target,images,album').then(response => {
          const graphPhoto = response as FbGraphPhoto;
          const photo = this.photoResolver.photoFromFBPhoto(graphPhoto);
          const promisePhoto = this.photos.get(photo.id);
          promisePhoto.photo = photo;
          this.photos.set(photo.id, promisePhoto);
          resolve(this.photos.get(photo.id).photo);
        }).catch(() => reject());
        this.photos.set(id, pP);
      } else {
        const pP = this.photos.get(id);
        if (pP.photo) {
          resolve(this.photos.get(id).photo);
        } else {
          pP.promise.then(() => resolve(this.photos.get(id).photo))
        }
      }
    });
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
          const photo = this.photoResolver.photoFromFBPhoto(fbPhoto);

          if (goodUserIds.includes(fbPhoto.from.id)) {
            photo.negotiable = true;
            negotiablePhotos.push(photo);
          }
          let p;
          if (this.photos.has(photo.id)) {
            p = this.photos.get(photo.id);
          } else {
            p = new PromisePhoto();
          }
          p.photo = photo;
          this.photos.set(photo.id, p);
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
    const nP = this.photos.get(photo.id);
    nP.photo = p;
    this.photos.set(photo.id, nP);
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
