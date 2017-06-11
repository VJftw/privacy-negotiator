import { Injectable } from '@angular/core';
import { Headers, Http } from '@angular/http';
import 'rxjs/add/operator/toPromise';
import { environment } from '../environments/environment';
import { WebSocketService } from './websocket.service';

@Injectable()
export class APIService {

  private headers = new Headers({'Content-Type': 'application/json'});
  private websocket: WebSocket;

  constructor(
    private http: Http,
    private webSocketService: WebSocketService
  ) {

  }

  public setAuthorization(authToken: string) {
    console.log('Authenticated with API:' + authToken);
    this.headers.set('Authorization', 'bearer ' + authToken);
    // connect to WebSocket
    this.websocket = new WebSocket(environment.apiEndpoint.replace('http', 'ws') + '/v1/ws?authToken=' + authToken);

    this.websocket.addEventListener('open', this.webSocketService.onOpen);
    this.websocket.addEventListener('message', this.webSocketService.onMessage);
    this.websocket.addEventListener('close', this.webSocketService.onClose);
    this.websocket.addEventListener('error', this.webSocketService.onError);

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
