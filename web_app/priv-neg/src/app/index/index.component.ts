import { Component } from '@angular/core';
import { Headers, Http } from '@angular/http';
import 'rxjs/add/operator/toPromise';
import { FacebookService, InitParams, LoginResponse } from 'ngx-facebook';

@Component({
  selector: 'index',
  templateUrl: './index.component.html',
})
export class IndexComponent {

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

  loginWithFacebook(): void {
    this.fb.login()
      .then((response: LoginResponse) => this.authenticateWithApi(response))
      .catch((error: any) => console.error(error));
  }

  authenticateWithApi(response: LoginResponse): void {
    this.http.post(
      "http://localhost/v1/auth",
      JSON.stringify(response.authResponse),
      {headers: new Headers({'Content-Type': 'application/json'})}
    ).toPromise()
    .then(res => console.log(res.json()))
    .catch(this.handleError)
    ;


  }

  private handleError(error: any): Promise<any> {
    console.error('An error occurred', error); // for demo purposes only
    return Promise.reject(error.message || error);
  }
}
