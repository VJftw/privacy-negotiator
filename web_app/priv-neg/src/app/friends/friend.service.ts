import { Injectable } from '@angular/core';
import { FacebookService } from 'ngx-facebook';
import {APIFriend, FBFriend, Friend} from './friend.model';
import {FBUser} from '../auth.service';
import {APIService} from '../api.service';
import {APIClique, Clique} from './clique.model';
import {Channel} from '../websocket.service';


@Injectable()
export class FriendService implements Channel {

  private friends: Map<string, Friend>;
  private cliques: Map<string, Clique>;
  protected offset: string;


  constructor(
    private fb: FacebookService,
    private apiService: APIService
  ) {
    this.friends = new Map();
    this.cliques = new Map();
  }

  public getFriends() {
    return Array.from(this.friends.values());
  }

  public getName(): string {
    return 'clique';
  }

  public onWebsocketMessage(data) {
    const apiClique = data as APIClique;
    if (!this.cliques.has(apiClique.id)) {
      const c = new Clique();
      c.name = apiClique.name
      this.cliques.set(apiClique.id, c);
    } else {
      // Merge cliques
      // this.cliques.get(apiClique)
    }
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
        const f = Friend.UpdateFromAPIFriend(this.friends.get(apiFriend.id), apiFriend);
        this.friends.set(f.id, f);

        existingFriendIds.push(f.id);
      }

      for (const friend of fbFriends) {
        if (!existingFriendIds.includes(friend.id)) {
          nonExistingFriends.push(friend);
        }
      }

    });

  }


}
