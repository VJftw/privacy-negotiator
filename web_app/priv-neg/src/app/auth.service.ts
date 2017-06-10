import { Injectable } from '@angular/core';
import { Headers, Http } from '@angular/http';
import 'rxjs/add/operator/toPromise';
import { FacebookService, InitParams, LoginResponse, LoginOptions } from 'ngx-facebook';
import { environment } from '../environments/environment';
import { Router, CanActivate } from '@angular/router';



@Injectable()
export class AuthService implements CanActivate {

  private authToken: string;
  private fbUser: FBUser;

  constructor(
    public fb: FacebookService,
    private http: Http,
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
    if (this.authToken) {
      return true;
    }
    console.log('Not logged in');
    return false;
  }

  public getUser(): FBUser {
    return this.fbUser;
  }

  public authenticate(): Promise<any> {
    const options: LoginOptions = {
      scope: 'user_friends,user_photos,user_posts'
    };

    return this.fb.login(options)
      .then((response: LoginResponse) => {
        this.fbUser = new FBUser();
        this.fbUser.shortAccessToken = response.authResponse.accessToken;
        this.fbUser.id = response.authResponse.userID;
        this.authenticateWithApi(response)
          .then(() => {
            this.router.navigate(['/photos']);
          });
      })
      .catch((error: any) => console.error(error))
    ;
  }

  public authenticateWithApi(response: LoginResponse): Promise<any> {
    return this.http.post(
      environment.apiEndpoint + '/v1/auth',
      JSON.stringify(response.authResponse),
      {headers: new Headers({'Content-Type': 'application/json'})}
    ).toPromise()
    .then(res => this.storeAuthToken(res.json()))
    .catch(this.handleError)
    ;
  }

  private storeAuthToken(response) {
    console.log('Logged in.');
    this.authToken = response.authToken;
  }

  private handleError(error: any): Promise<any> {
    console.error('An error occurred', error); // for demo purposes only
    return Promise.reject(error.message || error);
  }
}

export class FBUser {
  id: string;
  name: string;
  shortAccessToken: string;
}
