import { Injectable } from '@angular/core';
import { FacebookService } from 'ngx-facebook';


@Injectable()
export class CommunityService {

  constructor(
    private fb: FacebookService
  ) {}

  public getFriendsForUser(userId: string, offset = null): Promise<any> {
    let uri = '/' + userId + '/friends?fields=id,name,picture{url}';

    if (offset) {
      uri += '&after=' + offset;
    }

    return this.fb.api(uri)
      .catch(e => console.error(e))
    ;
  }

}