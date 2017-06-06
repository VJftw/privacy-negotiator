import { AuthService } from '../auth.service';
import { Component } from '@angular/core';
import { Router } from '@angular/router';

@Component({
  selector: 'app-index',
  templateUrl: './index.component.html',
  providers: [
    AuthService
  ]
})
export class IndexComponent {

  constructor(
    private authService: AuthService,
    private router: Router,
  ) {}

  loginWithFacebook(): void {
    this.authService.authenticate();
    this.router.navigate(['photos']);
  }
}
