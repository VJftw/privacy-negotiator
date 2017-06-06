import { AuthService } from '../auth.service';
import { Component } from '@angular/core';

@Component({
  selector: 'index',
  templateUrl: './index.component.html',
  providers: [
    AuthService
  ]
})
export class IndexComponent {

  constructor(
    private authService: AuthService
  ) {}

  loginWithFacebook(): void {
    this.authService.authenticate();
  }
}
