export class WebSocketService {

  private channels: Map<string, Channel>;


  public onOpen(ev: Event) {
    console.log('Connected to Websocket.');

  }

  public onMessage(ev: MessageEvent) {
    console.log('Received message from websocket' + ev.data);
  }
}

export interface Channel {
  onWebsocketMessage(MessageEvent);
}
