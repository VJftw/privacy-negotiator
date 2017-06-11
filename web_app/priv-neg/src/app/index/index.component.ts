import { AuthService } from '../auth.service';
import { Component } from '@angular/core';
import { LoginResponse } from 'ngx-facebook';

@Component({
  selector: 'app-index',
  templateUrl: './index.component.html',
})
export class IndexComponent {

  protected loading: boolean;

  constructor(
    private authService: AuthService,
  ) {
    this.loading = false;
  }

  loginWithFacebook(): void {
    this.loading = true;
    this.authService.authenticate();
  }
}
