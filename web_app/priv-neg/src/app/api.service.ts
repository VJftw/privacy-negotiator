import { Injectable } from '@angular/core';
import { Headers, Http } from '@angular/http';
import 'rxjs/add/operator/toPromise';
import { environment } from '../environments/environment';

@Injectable()
export class APIService {

  private headers = new Headers({'Content-Type': 'application/json'});

  constructor(
    private http: Http
  ) {}

  public setAuthorization(authToken: string) {
    console.log('Authenticated with API:' + authToken);
    this.headers.set('Authorization', 'bearer ' + authToken);
  }

  public isAuthenticated(): boolean {
    if (this.headers.has('Authorization')) {
      return true;
    }
    console.log('Not authenticated with API.');
    return false;
  }

  public post(resource: string, body: any): Promise<any> {
    return this.http.post(
      environment.apiEndpoint + resource,
      JSON.stringify(body),
      {headers: this.headers}
    ).toPromise()
    .catch(this.handleError);
  }

  public get(resource: string): Promise<any> {
    return this.http.get(
      environment.apiEndpoint + resource,
      {headers: this.headers}
    ).toPromise()
    .catch(this.handleError);
  }

  private handleError(error: any): Promise<any> {
    console.error('An error occurred', error);
    return Promise.reject(error.message || error);
  }

}
