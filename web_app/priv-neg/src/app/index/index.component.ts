import { AuthService } from '../auth.service';
import { Component } from '@angular/core';

@Component({
  selector: 'app-index',
  templateUrl: './index.component.html',
  // Social login styles from https://terrymooreii.github.io/materialize-social/
  styles: [
    `
      .btn.social > :first-child,
      .btn-large.social > :first-child,
      .col.social > :first-child {
        border-right: 1px solid rgba(0, 0, 0, 0.2);
      }
      .btn.social,
      .btn-large.social,
      .col.social {
        padding: 0 2rem 0 0;
      }
      .btn.social i,
      .btn-large.social i,
      .col.social i {
        padding: 0 1rem;
        margin-right: 1rem;
      }
      .btn.social-icon,
      .btn-large.social-icon,
      .col.social-icon {
        padding: 0;
      }
      .btn.social-icon i,
      .btn-large.social-icon i,
      .col.social-icon i {
        padding: 0 1rem;
        margin-right: 0;
      }
      .btn-large.social-icon {
        padding: 0 1rem;
      }
      .adn {
        background-color: #d87a68;
        color: #fff !important;
      }
      .adn i {
        color: #fff !important;
      }
      .adn:hover {
        background-color: #e29e91 !important;
      }
      .facebook {
        background-color: #3b5998;
        color: #fff !important;
      }
      .facebook i {
        color: #fff !important;
      }
      .facebook:hover {
        background-color: #4c70ba !important;
      }
    `
  ]
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
