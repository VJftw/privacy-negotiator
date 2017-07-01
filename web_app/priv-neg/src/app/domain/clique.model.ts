
import {CategorySelection} from './category.model';
import {User} from './user.model';

export class Clique {
  id: string;
  name: string;
  friends: Map<string, User>;
  editing = false;
  categories: CategorySelection[] = [];

  constructor() {
    this.name = 'Unnamed';
    this.friends = new Map();
  }

  public getFriends(): User[] {
    return Array.from(this.friends.values());
  }

  public removeFriend(userID: string) {
    this.friends.delete(userID);
  }
}

export class APIClique {
  id: string;
  name: string;
  users: string[];
  categories: string[] = [];

  public static FromClique(c: Clique): APIClique {
    const a = new APIClique();

    a.id = c.id;
    a.name = c.name;
    a.users = [];

    for (const u of c.getFriends()) {
      a.users.push(u.id);
    }

    a.categories = [];
    for (const cat of c.categories) {
      if (cat.isActive) {
        a.categories.push(cat.name);
      }
    }

    return a;
  }
}
