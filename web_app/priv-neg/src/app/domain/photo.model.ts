import {User} from './user.model';

export class FbGraphPhoto {
  id: string;
  created_time: string;
  from: User;
  images: FbPlatformImageSource[];

  album: FbGraphAlbum;
}

export class FbGraphAlbum {
  id: string;
  name: string;
  from: User;
}

export class FbPlatformImageSource {
  height: number;
  source: string;
  width: number;
}

export class Photo {
  id: string;
  createdTime: string;
  url: string;
  albumId: string;
  from: User;
  negotiable = false;
  pending = false;
  taggedUsers: User[] = [];
  categories: string[] = [];
  conflicts: Conflict[] = [];
  allowedUsers: User[] = [];
  blockedUsers: User[] = [];

  constructor(id: string) {
    this.id = id;
  }
}

export class Conflict {

  public static RESULT_ALLOWED = 1;
  public static RESULT_BLOCKED = -1;
  public static RESULT_INDETERMINATE = 0;

  id: string;
  target: User;
  parties: User[] = [];
  reasoning: Reason[] = [];
  result: string;

  constructor(id: string) {
    this.id = id;
  }
}

export class Reason {
  user: User;
  vote: number;
}

export class APIPhoto {
  id: string;
  taggedUsers: string[] = [];
  pending = false;
  categories: string[] = [];
  conflicts: APIConflict[];
  allowedUsers: string[];
  blockedUsers: string[];
}

export class APIConflict {
  id: string;
  target: string;
  parties: string[] = [];
  reasoning: APIReason[] = [];
  result: string; // allow, block, indeterminate
}

export class APIReason {
  id: string;
  vote: number;
}
