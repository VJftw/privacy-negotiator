import { Injectable } from '@angular/core';
import { FacebookService } from 'ngx-facebook';
import {APIFriend, FBFriend, Friend} from './friend.model';
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
    const c = new Clique();
    c.id = 'NA';
    c.name = 'NA';
    c.friends = new Map();
    this.cliques.set(c.id, c);
  }

  public getCliques(): Clique[] {
    return Array.from(this.cliques.values());
  }

  public getName(): string {
    return 'clique';
  }

  public onWebsocketMessage(data) {
    const apiClique = data as APIClique;
    if (!this.cliques.has(apiClique.id)) {
      const c = new Clique();
      c.name = apiClique.name;
      for (const userId of apiClique.users) {
        c.friends.set(userId, this.friends.get(userId));
      }
      this.cliques.set(apiClique.id, c);
    } else {
      // Merge cliques
      // this.cliques.get(apiClique)
    }
  }

  public updateCliquesFromAPI(): Promise<any> {
    return this.apiService.get(
      '/v1/cliques'
    ).then(response => {
      const apiCliques = response.json() as APIClique[];

      for (const apiClique of apiCliques) {
        if (this.cliques.has(apiClique.id)) {
          const clique = this.cliques.get(apiClique.id);
          clique.name = apiClique.name;
          this.cliques.set(clique.id, clique);
        } else {
          const clique = new Clique();
          clique.id = apiClique.id;
          if (apiClique.name === '') {
            clique.name = 'Unnamed';
          } else {
            clique.name = apiClique.name;
          }
          this.cliques.set(clique.id, clique);
        }
      }
    });
  }


  public updateFriends(offset = null): Promise<any> {
    let uri = '/me/friends?fields=id,name,picture{url}';

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
      for (const apiFriend of response.json() as APIFriend[]) {
        if (apiFriend.cliques.length < 1) {
          const clique = this.cliques.get('NA');
          clique.friends.set(apiFriend.id, this.friends.get(apiFriend.id));
          this.cliques.set('NA', clique);
        } else {
          for (const cliqueID of apiFriend.cliques) {
            if (this.cliques.has(cliqueID)) {
              const clique = this.cliques.get(cliqueID);
              clique.friends.set(apiFriend.id, this.friends.get(apiFriend.id));

              this.cliques.set(cliqueID, clique);
            } else {
              const clique = new Clique();
              clique.id = cliqueID;
              clique.name = 'Unnamed';
              clique.friends = new Map();
              clique.friends.set(apiFriend.id, this.friends.get(apiFriend.id));
              this.cliques.set(cliqueID, clique);
            }
          }
        }
      }
    });

  }


}
