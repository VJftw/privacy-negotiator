import {FbGraphUser} from '../auth.service';

export class User {
  id: string;
  firstName: string;
  lastName: string;
  name: string;
  picture: string;
  tieStrength: number;

  public static FromFBGraphUser(fbGraphUser: FbGraphUser): User {
    const u = new User(
      fbGraphUser.id,
      fbGraphUser.first_name,
      fbGraphUser.last_name,
      fbGraphUser.picture.data.url
    );

    return u
  }

  constructor(
    id: string,
    firstName: string,
    lastName: string,
    picture: string,
  ) {
    this.id = id;
    this.firstName = firstName;
    this.lastName = lastName;
    this.name = firstName + ' ' + lastName;
    this.picture = picture;
  }

}
