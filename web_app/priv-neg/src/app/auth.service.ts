import { Injectable } from '@angular/core';
import { Headers, Http } from '@angular/http';
import 'rxjs/add/operator/toPromise';
import { FacebookService, InitParams, LoginResponse } from 'ngx-facebook';
import { environment } from '../environments/environment';


@Injectable()
export class AuthService {

  private authToken: string;

  constructor(
    private fb: FacebookService,
    private http: Http
  ) {
    let initParams: InitParams = {
      appId: '219608771883029',
      xfbml: true,
      version: 'v2.8'
    };

    fb.init(initParams);
  }

  public isAuthenticated(): boolean {
    if (this.authToken) {
      return true;
    }

    return false;
  }

  public authenticate() {
    this.fb.login()
      .then((response: LoginResponse) => this.authenticateWithApi(response))
      .catch((error: any) => console.error(error))
    ;
  }

  private authenticateWithApi(response: LoginResponse): void {
    this.http.post(
      environment.apiEndpoint + "/v1/auth",
      JSON.stringify(response.authResponse),
      {headers: new Headers({'Content-Type': 'application/json'})}
    ).toPromise()
    .then(res => this.storeAuthToken(res))
    .catch(this.handleError)
    ;
  }

  private storeAuthToken(response) {
    this.authToken = response.authToken;
  }

  private handleError(error: any): Promise<any> {
    console.error('An error occurred', error); // for demo purposes only
    return Promise.reject(error.message || error);
  }
}
