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
  conflict: Conflict;
  allowedUsers: User[];
  blockedUsers: User[];

  public static fromFBPhoto(fp: FbGraphPhoto): Photo {
    const p = new Photo();

    p.id = fp.id;
    p.createdTime = fp.created_time;
    p.url = fp.images[0].source;
    p.albumId = fp.album.id;
    p.from = fp.from;
    return p;
  }

  public static fromAPIPhoto(ap: APIPhoto, p: Photo = new Photo()): Photo {
    p.id = ap.id;
    p.pending = ap.pending;
    p.categories = ap.categories;
    p.conflict = ap.conflict;

    return p;
  }
}

export class Conflict {
  id: string;
  targets: string[] = [];
  parties: string[] = [];
  resolved: boolean;
}

export class APIPhoto {
  id: string;
  taggedUsers: string[] = [];
  pending = false;
  categories: string[] = [];
  conflict: Conflict;
  allowedUsers: string[];
  blockedUsers: string[];
}

