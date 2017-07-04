import { Injectable } from '@angular/core';
import 'rxjs/add/operator/toPromise';
import {CanActivate, Router} from '@angular/router';

@Injectable()
export class SessionService implements CanActivate {

  private sessionUser: SessionUser;

  constructor(
    private router: Router,
  ) {}

  public setUser(u: SessionUser) {
    this.sessionUser = u;
  }

  public getUser(): SessionUser {
    return this.sessionUser;
  }

  public isAuthenticated(): boolean {
    return this.sessionUser != null
  }

  public canActivate(): boolean {
    if (this.isAuthenticated()) {
      return true;
    }
    this.router.navigate(['/']);
    return false;
  }

}

export class SessionUser {
  id: string;
  firstName: string;
  lastName: string;
  shortAccessToken: string;
  picture: string;
  coverPicture: string;

  constructor(id: string, firstName: string, lastName: string, shortAccessToken: string, picture: string) {
    this.id = id;
    this.firstName = firstName;
    this.lastName = lastName;
    this.shortAccessToken = shortAccessToken;
    this.picture = picture;
  }
}

