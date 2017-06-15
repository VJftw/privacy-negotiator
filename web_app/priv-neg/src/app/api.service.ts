import { Injectable } from '@angular/core';
import { Headers, Http } from '@angular/http';
import 'rxjs/add/operator/toPromise';
import { environment } from '../environments/environment';
import { WebSocketService } from './websocket.service';

function delay(ms: number) {
  return new Promise(resolve => setTimeout(resolve, ms));
}

@Injectable()
export class APIService {

  private headers = new Headers({'Content-Type': 'application/json'});
  private authorization: string;
  private websocket: WebSocket;

  constructor(
    private http: Http,
    private webSocketService: WebSocketService
  ) {}

  public setAuthorization(authToken: string) {
    console.log('Authenticated with API:' + authToken);
    this.headers.set('Authorization', 'bearer ' + authToken);
    this.authorization = authToken;
    this.connectToWebSocket();
  }

  private connectToWebSocket() {
    this.websocket = new WebSocket(environment.apiEndpoint.replace('http', 'ws') + '/v1/ws?authToken=' + this.authorization);

    this.websocket.addEventListener('open', (ev) => this.webSocketService.onOpen(ev));
    this.websocket.addEventListener('message', (ev) => this.webSocketService.onMessage(ev));
    this.websocket.addEventListener('close', (ev) => this.onClose(ev));
    this.websocket.addEventListener('error', (ev: ErrorEvent) => this.onError(ev));
  }

  private onClose(ev: Event) {
    console.log('Websocket closed.');
    console.log('Reconnecting...');
    delay(2000).then(() => this.connectToWebSocket());
  }

  private onError(ev: ErrorEvent) {
    console.log('Websocket error.');
    // console.log('Reconnecting...');
    // delay(2000).then(() => this.connectToWebSocket());
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

  public put(resource: string, body: any): Promise<any> {
    return this.http.put(
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
