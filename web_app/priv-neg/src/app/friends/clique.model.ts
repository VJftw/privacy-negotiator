import {Friend} from './friend.model';
import {CategorySelection} from '../photos/photo-detail.component';

export class Clique {
  id: string;
  name: string;
  friends: Map<string, Friend>;
  editing = false;
  categories: CategorySelection[] = [];

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
  categories: string[] = [];
}
