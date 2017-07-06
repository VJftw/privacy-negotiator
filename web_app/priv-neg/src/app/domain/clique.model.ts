
import {ContextSelection} from './context.model';
import {User} from './user.model';

export class Clique {
  id: string;
  name: string;
  friends: Map<string, User>;
  editing = false;
  contexts: ContextSelection[] = [];

  constructor(id: string, name = 'Unnamed') {
    this.id = id;
    this.name = name;
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
    for (const cat of c.contexts) {
      if (cat.isActive) {
        a.categories.push(cat.name);
      }
    }

    return a;
  }
}
