import { Injectable } from '@angular/core';
import { FacebookService } from 'ngx-facebook';
import { Photo, FBPhoto } from './photo.model';
import { FBUser } from '../auth.service';


@Injectable()
export class PhotoService {

  protected photos: Photo[];
  protected offset: string;

  constructor(
    private fb: FacebookService
  ) {
    this.photos = [];
  }

  public getPhotos(): Photo[] {
    return this.photos;
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
        for (const fbPhoto of fbPhotos) {
          if (fbPhoto.album) {
            const photo = Photo.fromFBPhoto(fbPhoto);
            this.photos.push(photo);
          }
        }
      })
      .catch(e => console.error(e))
    ;
  }

  public getTaggedPhotosForUser(userId: string, offset = null): Promise<any> {
    let uri = '/' + userId + '/photos?fields=id,created_time,from,target,images,album';

    if (offset) {
      uri += '&after=' + offset;
    }

    return this.fb.api(uri)
      .catch(e => console.error(e))
    ;
  }

  public getTaggedPhotosForUsera(userId: string): Promise<Photo[]> {
    let uri = '/' + userId + '/photos?fields=id,created_time,from,target,images,album';

    if (this.offset) {
      uri += '&after=' + this.offset;
    }

    return this.fb.api(uri)
      .then(response => {
        if (response.paging) {
          this.offset = response.paging.cursors.after;
        }
        return response.data as Photo[];
      })
      .catch(e => console.error(e))
    ;
  }

  public getPrivacyForAnAlbum(albumId: string): Promise<any> {
    const uri = '/' + albumId + '?fields=id,privacy';

    return this.fb.api(uri)
      .catch(e => console.error(e))
    ;
  }
}
