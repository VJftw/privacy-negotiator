import { Injectable } from '@angular/core';
import 'rxjs/add/operator/toPromise';
import { FacebookService, InitParams, LoginResponse, LoginOptions } from 'ngx-facebook';
import { environment } from '../environments/environment';
import { Router, CanActivate } from '@angular/router';
import { APIService } from './api.service';
import {CategoryService} from './categories/category.service';
import {FriendService} from './friends/friend.service';
import {PhotoService} from './photos/photo.service';

@Injectable()
export class AuthService implements CanActivate {

  private fbUser: SessionUser;

  constructor(
    public fb: FacebookService,
    private apiService: APIService,
    private categoryService: CategoryService,
    private photoService: PhotoService,
    private friendSerivce: FriendService,
    private router: Router,

  ) {
    const initParams: InitParams = {
      appId: environment.fbAppId,
      xfbml: true,
      version: 'v2.9'
    };

    fb.init(initParams);
  }

  public canActivate(): boolean {
    if (this.isAuthenticated()) {
      return true;
    }
    this.router.navigate(['']);
    return false;
  }

  public isAuthenticated(): boolean {
    if (this.fbUser) {
      return this.apiService.isAuthenticated();
    }
    console.log('Not authenticated with Facebook.');
    return false;
  }

  public getUser(): SessionUser {
    return this.fbUser;
  }

  public authenticate(): Promise<any> {
    const options: LoginOptions = {
      scope:
      'user_friends,' +
      'user_photos,' +
      'user_posts,' +
      'user_education_history,' +
      'user_hometown,' +
      'user_likes,' +
      'user_location,' +
      'user_relationship_details,' +
      'user_relationships,' +
      'user_religion_politics,' +
      'user_work_history,' +
      'user_events,' +
      'user_managed_groups'
    };

    return this.fb.login(options)
      .then((response: LoginResponse) => {
        console.log('Authenticated with Facebook. Getting /me');
        this.fb.api('/me?fields=id,name,picture{url},cover').then(res => {
          const user = res as FbGraphUser;
          console.log(user);
          this.fbUser = new SessionUser(
            response.authResponse.userID,
            user.name,
            response.authResponse.accessToken,
            user.picture.data.url
          );
          console.log(user);
          if (user.cover && user.cover.source) {
            this.fbUser.coverPicture = user.cover.source;
          }
          this.authenticateWithApi(response)
            .then(() => {
              this.categoryService.updateCategories().then(() => {
                this.apiService.webSocketService.addChannel(this.photoService);
                this.apiService.webSocketService.addChannel(this.friendSerivce);
                this.router.navigate(['/photos']);
              });
            });
        });
      })
      .catch((error: any) => console.error(error))
    ;
  }

  private authenticateWithApi(loginResponse: LoginResponse): Promise<any> {
    return this.apiService.post(
      '/v1/auth',
      loginResponse.authResponse
    ).then(
      authResponse => this.apiService.setAuthorization(authResponse.json().authToken)
    );
  }

  private handleError(error: any): Promise<any> {
    console.error('An error occurred', error); // for demo purposes only
    return Promise.reject(error.message || error);
  }
}

export class SessionUser {
  id: string;
  name: string;
  shortAccessToken: string;
  picture: string;
  coverPicture: string;

  constructor(id: string, name: string, shortAccessToken: string, picture: string) {
    this.id = id;
    this.name = name;
    this.shortAccessToken = shortAccessToken;
    this.picture = picture;
  }
}

export class FbGraphUser {
  id: string;
  name: string;
  picture: FbGraphUserPicture;
  cover: FbGraphCover;
}

export class FbGraphUserPicture {
  data: FbGraphUserPictureData;
}


export class FbGraphUserPictureData {
  url: string;
}

export class FbGraphCover {
  id: string;
  source: string;
}
