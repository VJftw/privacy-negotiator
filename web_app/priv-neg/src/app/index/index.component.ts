import { AuthService } from '../auth.service';
import { Component } from '@angular/core';
import { LoginResponse } from 'ngx-facebook';

@Component({
  selector: 'app-index',
  templateUrl: './index.component.html',
})
export class IndexComponent {

  constructor(
    private authService: AuthService,
  ) {}

  loginWithFacebook(): void {
    this.authService.authenticate();
  }
}
