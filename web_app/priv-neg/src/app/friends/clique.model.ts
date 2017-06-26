import {Friend} from './friend.model';

export class Clique {
  id: string;
  name: string;
  friends: Map<string, Friend>;

  constructor() {
    this.friends = new Map();
  }

  public getFriends(): Friend[] {
    return Array.from(this.friends.values());
  }
}

export class APIClique {
  id: string;
  name: string;
  users: string[];
}
