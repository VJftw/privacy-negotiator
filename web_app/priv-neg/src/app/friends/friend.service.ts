import { Injectable } from '@angular/core';
import { FacebookService } from 'ngx-facebook';
import {APIFriend, FBFriend, Friend} from './friend.model';
import {FBUser} from '../auth.service';
import {APIService} from '../api.service';


@Injectable()
export class FriendService {

  private friends: Map<string, Friend>;
  protected offset: string;


  constructor(
    private fb: FacebookService,
    private apiService: APIService
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
    const friendIds = [];
    for (const fbFriend of fbFriends) {
      const friend = Friend.FromFBFriend(fbFriend);
      friendIds.push(fbFriend.id);

      this.friends.set(friend.id, friend);
    }

    this.apiService.get(
      '/v1/friends?ids=' + JSON.stringify(friendIds)
    ).then(response => {
      const existingFriends = response.json() as APIFriend[];
      const existingFriendIds = [];
      const nonExistingFriends = [];

      for (const apiFriend of existingFriends) {
        const f = Friend.UpdateFromAPIFriend(this.friends.get(apiFriend.to), apiFriend);
        this.friends.set(f.id, f);

        existingFriendIds.push(f.id);
      }

      for (const friend of fbFriends) {
        if (!existingFriendIds.includes(friend.id)) {
          nonExistingFriends.push(friend);
        }
      }

      this.addNewFriendsToAPI(nonExistingFriends);

    });

  }

  private addNewFriendsToAPI(newFriends: FBFriend[]) {

    for (const friend of newFriends) {
      this.apiService.post(
        '/v1/friends',
        friend
      ).then(response => console.log(response));
    }
  }


}
