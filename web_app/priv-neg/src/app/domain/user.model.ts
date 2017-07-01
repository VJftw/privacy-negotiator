import {FbGraphUser} from '../auth.service';

export class User {
  id: string;
  name: string;
  picture: string;
  tieStrength: number;

  public static FromFBGraphUser(fbGraphUser: FbGraphUser): User {
    const u = new User(
      fbGraphUser.id,
      fbGraphUser.name,
      fbGraphUser.picture.data.url
    );

    return u
  }

  constructor(
    id: string,
    name: string,
    picture: string,
  ) {
    this.id = id;
    this.name = name;
    this.picture = picture;
  }

}
