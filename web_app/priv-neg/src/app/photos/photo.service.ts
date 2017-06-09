import { Injectable } from '@angular/core';
import { FacebookService } from 'ngx-facebook';
import { Photo } from './photo.model';


@Injectable()
export class PhotoService {

  protected offset: string;

  constructor(
    private fb: FacebookService
  ) {}

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
