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

  }

}
