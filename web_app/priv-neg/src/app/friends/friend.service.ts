import { Injectable } from '@angular/core';
import { FacebookService } from 'ngx-facebook';
import {FBFriend, Friend} from "./friend.model";
import {FBUser} from "../auth.service";


@Injectable()
export class FriendService {

  private friends: Map<string, Friend>;
  protected offset: string;


  constructor(
    private fb: FacebookService
  ) {
    this.friends = new Map();
  }

  public getFriends() {
    return Array.from(this.friends.values());
  }

  public updateFriendsForUser(fbUser: FBUser, offset = null): Promise<any> {
    let uri = '/' + fbUser.id + '/friends?fields=id,name,picture{url}';

    if (offset) {
      uri += '&after=' + offset;
    }

    return this.fb.api(uri)
      .then(response => {
        if (response.paging) {
          this.offset = response.paging.cursors.after;
        }
        const fbFriends = response.data as FBFriend[];
        console.log(fbFriends);
        this.processFriends(fbFriends);
      })
      .catch(e => console.error(e))
    ;
  }

  private processFriends(fbFriends: FBFriend[]) {
    for (const fbFriend of fbFriends) {
      const friend = Friend.FromFBFriend(fbFriend);
      console.log(friend);

      this.friends.set(friend.id, friend)
    }
  }

}
