import {AfterViewChecked, Component, OnDestroy, OnInit} from '@angular/core';
import { FriendService } from './friend.service';
import {NavigationEnd, Router} from '@angular/router';

declare var Materialize: any;

@Component({
  selector: 'app-friends',
  templateUrl: './friends.component.html'
})
export class FriendsComponent implements OnInit, AfterViewChecked, OnDestroy {

  public lock = false;
  private _subscription;

  constructor(
    public friendService: FriendService,
    private router: Router,
  ) {
  }

  ngOnInit() {
    this.friendService.updateCliquesFromAPI().then(() => this.updateFriends());
    this._subscription = this.router.events.subscribe((e: NavigationEnd) => {
      if (e instanceof  NavigationEnd && e.url === '/friends') {
        console.log(e);
        this.friendService.updateCliquesFromAPI().then(() => this.updateFriends());
      }
    });
  }

  updateFriends() {
    if (!this.lock) {
      this.lock = true;
      this.friendService.updateFriends()
      .then(() => this.lock = false)
      ;
    }
  }

  toggleEdit(cliqueId: string) {
    const clique = this.friendService.getCliqueById(cliqueId);
    if (clique.editing === false) {
      clique.editing = true;
    } else {
      clique.editing = false;
      this.friendService.updateClique(cliqueId, clique);
    }
    console.log(clique);
  }

  ngAfterViewChecked() {
    Materialize.updateTextFields();
  }

  ngOnDestroy() {
    this._subscription.unsubscribe();
  }

}
