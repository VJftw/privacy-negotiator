import { Component } from '@angular/core';
import { FacebookService, InitParams, LoginResponse } from 'ngx-facebook';

@Component({
  selector: 'index',
  templateUrl: './index.component.html',
})
export class IndexComponent {

  constructor(private fb: FacebookService) {

    let initParams: InitParams = {
      appId: '219608771883029',
      xfbml: true,
      version: 'v2.8'
    };

    fb.init(initParams);
  }

  loginWithFacebook(): void {

    this.fb.login()
      .then((response: LoginResponse) => console.log(response))
      .catch((error: any) => console.error(error));

  }
}
