import { Injectable } from '@angular/core';
import { FacebookService } from 'ngx-facebook';


@Injectable()
export class PhotoService {

  constructor(
    private fb: FacebookService
  ) {}

  public getTaggedPhotosForUser(userId: string, offset = null): Promise<any> {
    let uri = '/' + userId + '/photos?fields=id,created_time,from,source,album';

    if (offset) {
      uri += '&after=' + offset;
    }

    return this.fb.api(uri)
      .catch(e => console.error(e))
    ;
  }
}
