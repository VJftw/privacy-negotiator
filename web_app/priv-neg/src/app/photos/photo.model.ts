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
  uploader: FBUser;

  public static fromFBPhoto(fp: FBPhoto): Photo {
    const p = new Photo();

    p.id = fp.id;
    p.createdTime = fp.created_time;
    p.url = fp.images[0].source;
    p.albumId = fp.album.id;
    p.uploader = fp.from;

    return p;
  }
}
