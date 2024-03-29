import {Injectable} from '@angular/core';
import { FacebookService } from 'ngx-facebook';
import {APIFriend} from '../domain/friend.model';
import {APIService} from '../api.service';
import {APIClique, Clique} from '../domain/clique.model';
import {Channel} from '../websocket.service';
import {ContextService} from '../contexts/context.service';
import {ContextSelection} from '../domain/context.model';
import {FbGraphUser} from '../auth.service';
import {User} from '../domain/user.model';

class PromiseUser {
  promise: Promise<any>;
  user: User;
}

@Injectable()
export class FriendService implements Channel {

  private friends: Map<string, PromiseUser>;
  private cliques: Map<string, Clique>;
  protected offset: string;


  constructor(
    private fb: FacebookService,
    private apiService: APIService,
    private contextService: ContextService,
  ) {
    this.friends = new Map();
    this.cliques = new Map();
  }

  private resetCliques() {
    this.cliques.clear();
    const c = new Clique('NA', 'Not Grouped');
    c.friends = new Map();
    this.cliques.set(c.id, c);
  }

  public getUserById(id: string): Promise<User> {
    return new Promise((resolve, reject) => {
      if (!this.friends.has(id)) {
        const pU = new PromiseUser();
        pU.promise = this.fb.api('/' + id + '?fields=id,first_name,last_name,picture.type(large){url}').then(response => {
          const friend = response as FbGraphUser;
          const user = User.FromFBGraphUser(friend);
          const promiseUser = this.friends.get(user.id);
          promiseUser.user = user;
          this.friends.set(user.id, promiseUser);
          resolve(this.friends.get(user.id).user);
        });
        this.friends.set(id, pU);
      } else {
        const pU = this.friends.get(id);
        if (pU.user) {
          resolve(this.friends.get(id).user);
        } else {
          pU.promise.then(() => resolve(this.friends.get(id).user));
        }
      }
    });
  }

  public getCliques(): Clique[] {
    return Array.from(this.cliques.values());
  }

  public getName(): string {
    return 'clique';
  }

  public onWebSocketMessage(data) {
    this.updateCliquesFromAPI().then(() => this.updateFriends());
  }

  public updateCliquesFromAPI(): Promise<any> {
    return this.apiService.get(
      '/v1/cliques'
    ).then(response => {
      const apiCliques = response.json() as APIClique[];
      console.log(apiCliques);

      this.resetCliques();

      for (const apiClique of apiCliques) {
        let clique: Clique;
        if (this.cliques.has(apiClique.id)) {
          clique = this.cliques.get(apiClique.id);
        } else {
          clique = new Clique(apiClique.id);
        }
        if (apiClique.name.length < 1) {
          clique.name = 'Unnamed';
        } else {
          clique.name = apiClique.name;
        }
        clique.contexts = [];
        for (const cat of this.contextService.getContexts()) {
          let context;
          if (apiClique.categories.includes(cat.name)) {
            context = new ContextSelection(cat.name, true);
          } else {
            context = new ContextSelection(cat.name, false);
          }
          clique.contexts.push(context);
        }
        this.cliques.set(clique.id, clique);
      }
    });
  }

  public getCliqueById(id: string): Clique {
    if (this.cliques.has(id)) {
      return this.cliques.get(id);
    }
    return null;
  }

  public updateClique(id: string, clique: Clique) {
    this.cliques.set(id, clique);
    const uri = '/v1/cliques/' + id;

    this.apiService.put(uri, APIClique.FromClique(clique))
      .then(res => console.log(res))
    ;
  }


  public updateFriends(offset = null): Promise<any> {
    let uri = '/me/friends?fields=id,first_name,last_name,picture.type(large){url}&limit=500';

    if (offset) {
      uri += '&after=' + offset;
    }

    return this.fb.api(uri)
      .then(response => {
        if (response.paging) {
          this.offset = response.paging.cursors.after;
        }
        const fbFriends = response.data as FbGraphUser[];
        this.processFriends(fbFriends);
      })
      .catch(e => console.error(e))
    ;
  }

  private processFriends(fbFriends: FbGraphUser[]) {
    const friendIds = [];
    for (const fbFriend of fbFriends) {
      const friend = User.FromFBGraphUser(fbFriend);
      friendIds.push(fbFriend.id);

      const promiseUser = new PromiseUser();
      promiseUser.user = friend;

      this.friends.set(friend.id, promiseUser);
    }

    this.apiService.get(
      '/v1/friends?ids=' + JSON.stringify(friendIds)
    ).then(response => {
      for (const apiFriend of response.json() as APIFriend[]) {
        const friend = this.friends.get(apiFriend.id);
        friend.user = User.UpdateFromAPIFriend(friend.user, apiFriend);
        console.log(friend.user);
        this.friends.set(friend.user.id, friend);
        if (apiFriend.cliques.length <= 0) {
          const clique = this.cliques.get('NA');
          clique.friends.set(apiFriend.id, this.friends.get(apiFriend.id).user);
          this.cliques.set('NA', clique);
        } else {
          for (const cliqueID of apiFriend.cliques) {
            if (this.cliques.has(cliqueID)) {
              const clique = this.cliques.get(cliqueID);
              clique.friends.set(apiFriend.id, this.friends.get(apiFriend.id).user);
              this.cliques.get('NA').removeFriend(apiFriend.id);
              this.cliques.set(cliqueID, clique);
            } else if (cliqueID.length > 0) {
              const clique = new Clique(cliqueID);
              for (const cat of this.contextService.getContexts()) {
                clique.contexts.push(new ContextSelection(cat.name, false));
              }
              clique.friends = new Map();
              clique.friends.set(apiFriend.id, this.friends.get(apiFriend.id).user);
              this.cliques.set(clique.id, clique);
            }
          }
        }
      }
    });

  }
}
