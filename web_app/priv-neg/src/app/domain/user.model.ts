import {FbGraphUser} from '../auth.service';
import {APIFriend, TieStrengthDetail} from './friend.model';

export class User {
  id: string;
  firstName: string;
  lastName: string;
  name: string;
  picture: string;
  tieStrength: number;
  tieStrengthDetail: TieStrengthDetail[];

  public static FromFBGraphUser(fbGraphUser: FbGraphUser): User {
    const u = new User(
      fbGraphUser.id,
      fbGraphUser.first_name,
      fbGraphUser.last_name,
      fbGraphUser.picture.data.url
    );

    return u
  }

  public static UpdateFromAPIFriend(user: User, apiFriend: APIFriend): User {
    user.tieStrength = apiFriend.tieStrength;
    user.tieStrengthDetail = [];

    if (apiFriend.tieStrengthDetails) {

      for (const k in TieStrengthDetail.VALS) {
        if (apiFriend.tieStrengthDetails.hasOwnProperty(k)) {
          user.tieStrengthDetail.push(new TieStrengthDetail(
            k,
            TieStrengthDetail.VALS[k].humanKey,
            TieStrengthDetail.VALS[k].bool,
            apiFriend.tieStrengthDetails[k]
          ));
        }

      }
    }


    return user;
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
