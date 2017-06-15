import { Injectable } from '@angular/core';
import { FacebookService } from 'ngx-facebook';
import { Photo, FBPhoto, APIPhoto } from './photo.model';
import { FBUser } from '../auth.service';
import { APIService } from '../api.service';
import { WebSocketService, Channel } from '../websocket.service';


@Injectable()
export class PhotoService implements Channel {

  protected photos: Map<string, Photo>;
  protected offset: string;

  constructor(
    private fb: FacebookService,
    private apiService: APIService,
    private websocketService: WebSocketService
  ) {
    this.photos = new Map();
    this.websocketService.addChannel(this);
  }

  public getName(): string {
    return 'photo';
  }

  public onWebsocketMessage(data) {
    const apiPhoto = data as APIPhoto;
    const photo = this.getPhotoById(apiPhoto.id);
    this.photos.set(apiPhoto.id, Photo.fromAPIPhoto(apiPhoto, photo));
  }

  public getPhotos(): Photo[] {
    return Array.from(this.photos.values());
  }

  public updatePhoto(photo: Photo) {
    this.photos.set(photo.id, photo);
    // send PUT request to API
    this.apiService.put(
      '/v1/photos/' + photo.id,
      photo
    ).then(response => {
      console.log(response);
    });
  }

  public getPhotoById(photoId: string): Photo {
    return this.photos.get(photoId);
  }

  public updateTaggedPhotosForUser(fbUser: FBUser): Promise<Photo[]> {
    let uri = '/' + fbUser.id + '/photos?fields=id,created_time,from,target,images,album';

    if (this.offset) {
      uri += '&after=' + this.offset;
    }

    return this.fb.api(uri)
      .then(response => {
        if (response.paging) {
          this.offset = response.paging.cursors.after;
        }
        const fbPhotos = response.data as FBPhoto[];
        this.processFBPhotos(fbPhotos);
      })
      .catch(e => console.error(e))
    ;
  }

  private processFBPhotos(fbPhotos: FBPhoto[]) {

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
          const photo = Photo.fromFBPhoto(fbPhoto);

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
      console.log(foundPhotos);

      for (const photo of negotiablePhotos) {
        let f = null;
        for (const foundPhoto of foundPhotos) {
          if (foundPhoto.id === photo.id) {
            f = foundPhoto;
            break;
          }
        }
        if (f) {
          this.photos.set(photo.id, Photo.fromAPIPhoto(f, photo));
        } else {
          // POST new photo
          this.savePhoto(photo);
        }
      }
    });
  }

  public savePhoto(photo: Photo) {
    this.apiService.post(
      '/v1/photos',
      photo
    ).then(response => {
      console.log(response);
    });
  }

}
