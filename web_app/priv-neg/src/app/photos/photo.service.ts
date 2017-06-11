import { Injectable } from '@angular/core';
import { FacebookService } from 'ngx-facebook';
import { Photo, FBPhoto } from './photo.model';
import { FBUser } from '../auth.service';
import { APIService } from '../api.service';


@Injectable()
export class PhotoService {

  protected photos: Map<string, Photo>;
  protected offset: string;

  constructor(
    private fb: FacebookService,
    private apiService: APIService
  ) {
    this.photos = new Map();
  }

  public getPhotos(): Photo[] {
    return Array.from(this.photos.values());
  }

  public updatePhoto(photo: Photo) {
    this.photos.set(photo.id, photo);
    // send PUT request to API
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

    const userIds = [];

    // remove photos with no album id
    const fbPhotosWAlbum: FBPhoto[] = [];
    for (const fbPhoto of fbPhotos) {
      if (fbPhoto.album) {
        fbPhotosWAlbum.push(fbPhoto);
      }
    }

    for (const fbPhoto of fbPhotosWAlbum) {
      if (!userIds.includes(fbPhoto.from.id)) {
        userIds.push(fbPhoto.from.id);
      }
    }
    console.log(JSON.stringify(userIds));

    const photos: Photo[] = [];

    this.apiService.get(
      '/v1/users?ids=' + JSON.stringify(userIds)
    ).then(response => {
      const goodUserIds = response.json();
      console.log(goodUserIds);

      for (const fbPhoto of fbPhotosWAlbum) {
        const photo = Photo.fromFBPhoto(fbPhoto);
        if (goodUserIds.includes(fbPhoto.from.id)) {
          console.log(fbPhoto.from.id, fbPhoto.from.name);
          photo.negotiable = true;
        }
        photos.push(photo);
        this.photos.set(photo.id, photo);
      }
      console.log(this.photos);
    });
  }
}
