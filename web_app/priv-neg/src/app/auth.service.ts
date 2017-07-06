import { Injectable } from '@angular/core';
import 'rxjs/add/operator/toPromise';
import { FacebookService, InitParams, LoginResponse, LoginOptions } from 'ngx-facebook';
import { environment } from '../environments/environment';
import { Router } from '@angular/router';
import { APIService } from './api.service';
import {CategoryService} from './categories/category.service';
import {FriendService} from './friends/friend.service';
import {PhotoService} from './photos/photo.service';
import {SessionService, SessionUser} from './session.service';

@Injectable()
export class AuthService {

  constructor(
    public fb: FacebookService,
    private apiService: APIService,
    private sessionService: SessionService,
    private categoryService: CategoryService,
    private photoService: PhotoService,
    private friendService: FriendService,
    private router: Router,

  ) {
    const initParams: InitParams = {
      appId: environment.fbAppId,
      xfbml: true,
      version: 'v2.9'
    };

    fb.init(initParams);
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
        this.fb.api('/me?fields=id,first_name,last_name,picture{url},cover').then(res => {
          const user = res as FbGraphUser;
          console.log(user);
          const fbUser = new SessionUser(
            response.authResponse.userID,
            user.first_name,
            user.last_name,
            response.authResponse.accessToken,
            user.picture.data.url
          );
          console.log(user);
          if (user.cover && user.cover.source) {
            fbUser.coverPicture = user.cover.source;
          }
          this.sessionService.setUser(fbUser);
          this.authenticateWithApi(response)
            .then(() => {
              this.categoryService.updateCategoriesFromAPI().then(() => {
                this.apiService.webSocketService.addChannel(this.photoService);
                this.apiService.webSocketService.addChannel(this.friendService);
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

}

export class FbGraphUser {
  id: string;
  first_name: string;
  last_name: string;
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
