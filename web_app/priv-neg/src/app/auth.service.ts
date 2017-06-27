import { Injectable } from '@angular/core';
import 'rxjs/add/operator/toPromise';
import { FacebookService, InitParams, LoginResponse, LoginOptions } from 'ngx-facebook';
import { environment } from '../environments/environment';
import { Router, CanActivate } from '@angular/router';
import { APIService } from './api.service';
import {CategoryService} from './categories/category.service';

@Injectable()
export class AuthService implements CanActivate {

  private fbUser: SessionUser;

  constructor(
    public fb: FacebookService,
    private apiService: APIService,
    private categoryService: CategoryService,
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
      scope: 'user_friends,user_photos,user_posts'
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
              this.categoryService.updateCategories();
              this.router.navigate(['/photos']);
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

class FbGraphUser {
  id: string;
  name: string;
  picture: FbGraphUserPicture;
  cover: FbGraphCover;
}

class FbGraphUserPicture {
  data: FbGraphUserPictureData;
}


class FbGraphUserPictureData {
  url: string;
}

class FbGraphCover {
  id: string;
  source: string;
}
