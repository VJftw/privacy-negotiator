import { Injectable } from '@angular/core';

import { PhotoService } from './photos/photo.service';
import { Photo } from './photos/photo.model';
import {FriendService} from './friends/friend.service';

@Injectable()
export class WebSocketService {

  private channels: Map<string, Channel>;

  constructor() {
    this.channels = new Map();
  }

  public addChannel(c: Channel) {
    this.channels.set(
      c.getName(),
      c
    );
    console.log(this.channels);
  }

  public onOpen(ev: Event) {
    console.log('Connected to Websocket.');

  }

  public onMessage(ev: MessageEvent) {
    console.log('Received message from websocket' + ev.data);

    const wsMessage = JSON.parse(ev.data) as WSMessage;

    if (this.channels.has(wsMessage.type)) {
      const channel = this.channels.get(wsMessage.type);
      channel.onWebsocketMessage(wsMessage.data);
    } else {
      console.error('No channel found for: ' + wsMessage.type);
    }

  }
}

export interface Channel {
  getName(): string;
  onWebsocketMessage(data);
}

class WSMessage {
  type: string;
  data: any;
}
