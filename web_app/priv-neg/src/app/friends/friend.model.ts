import {FBPlatformImageSource} from "../photos/photo.model";

export class Friend {
  id: string;
  picture: string;
  name: string;

  public static FromFBFriend(fbFriend: FBFriend): Friend {
    let f = new Friend();

    f.id = fbFriend.id;
    f.picture = fbFriend.picture.data.url;
    f.name = fbFriend.name;

    return f;
  }
}

export class FBFriend {
  id: string;
  name: string;
  picture;
}

