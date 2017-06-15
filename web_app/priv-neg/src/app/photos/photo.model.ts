import { FBUser } from '../auth.service';

export class FBPhoto {
  id: string;
  created_time: string;
  from: FBUser;
  images: FBPlatformImageSource[];

  album: FBAlbum;
}

export class FBAlbum {
  id: string;
  name: string;
  from: FBUser;
}

export class FBPlatformImageSource {
  height: number;
  source: string;
  width: number;
}

export class Photo {
  id: string;
  createdTime: string;
  url: string;
  albumId: string;
  from: FBUser;
  negotiable = false;
  pending = false;
  taggedUsers: FBUser[] = [];
  categories: string[] = [];

  public static fromFBPhoto(fp: FBPhoto): Photo {
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

    for (const taggedUser of ap.taggedUsers) {
      const fbUser = new FBUser();
      fbUser.id = taggedUser;
      p.taggedUsers.push(fbUser);
    }

    p.categories = ap.categories;

    return p;
  }
}

export class APIPhoto {
  id: string;
  taggedUsers: string[] = [];
  pending = false;
  categories: string[] = [];
}
