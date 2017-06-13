import { Component } from '@angular/core';
import { environment } from '../environments/environment';
import { AuthService } from './auth.service';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html'
})
export class AppComponent {
  version = environment.version;

  constructor(
    private authService: AuthService
  ) {}

  protected isLoggedIn(): boolean {
    return this.authService.isAuthenticated();
  }
}
