import {Friend} from './friend.model';

export class Clique {
  id: string;
  name: string;
  friends: Friend[];
}

export class APIClique {
  id: string;
  name: string;
  users: string[];
}
