import { AuthService } from '../auth.service';
import { Component, OnInit } from '@angular/core';
import { FacebookService, InitParams, LoginResponse, LoginOptions } from 'ngx-facebook';
import { CommunityService } from './community.service';


@Component({
  selector: 'app-communities',
  templateUrl: './communities.component.html',
  providers: [
    CommunityService
  ]
})
export class CommunitiesComponent implements OnInit {

  private offset: string;

  constructor(
    private authService: AuthService,
    private communityService: CommunityService
  ) {}

  ngOnInit() {
    this.getFriends();
  }

  getFriends() {
    if (this.offset) {
      this.communityService.getFriendsForUser(
        this.authService.userId,
        this.offset
      ).then(res => {
        console.log(res);
      });
    } else {
      this.communityService.getFriendsForUser(
        this.authService.userId
      ).then(res => {
        console.log(res);
      });
    }
  }

}
