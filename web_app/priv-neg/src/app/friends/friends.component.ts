import { AuthService } from '../auth.service';
import { Component, OnInit } from '@angular/core';
import { FriendService } from './friend.service';
import {Friend} from './friend.model';


@Component({
  selector: 'app-friends',
  templateUrl: './friends.component.html'
})
export class FriendsComponent implements OnInit {

  private lock = false;

  constructor(
    private authService: AuthService,
    private friendService: FriendService
  ) {}

  ngOnInit() {
    this.updateFriends();
  }

  getFriends(): Friend[] {
    return this.friendService.getFriends();
  }

  updateFriends() {
    if (!this.lock) {
      this.lock = true;
      this.friendService.updateFriendsForUser(this.authService.getUser())
      .then(() => this.lock = false)
      ;
    }
  }

}
