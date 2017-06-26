import { AuthService } from '../auth.service';
import { Component, OnInit } from '@angular/core';
import { FriendService } from './friend.service';


@Component({
  selector: 'app-friends',
  templateUrl: './friends.component.html'
})
export class FriendsComponent implements OnInit {

  public lock = false;

  constructor(
    private authService: AuthService,
    public friendService: FriendService
  ) {}

  ngOnInit() {
    this.friendService.updateCliquesFromAPI().then(() => this.updateFriends());
  }

  updateFriends() {
    if (!this.lock) {
      this.lock = true;
      this.friendService.updateFriends()
      .then(() => this.lock = false)
      ;
    }
  }

}