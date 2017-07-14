import { AuthService } from '../auth.service';
import { Component } from '@angular/core';

@Component({
  selector: 'app-index',
  templateUrl: './index.component.html',
})
export class IndexComponent {

  public loading: boolean;
  public learnMore = false;
  public logInError = '';

  constructor(
    private authService: AuthService,
  ) {
    this.loading = false;
  }

  loginWithFacebook(): void {
    this.loading = true;
    this.authService.authenticate()
      .catch((error: any) => {
        this.logInError = error;
        this.loading = false;
      })
    ;
  }

  public toggleLearnMore() {
    this.learnMore = !this.learnMore;
  }
}
