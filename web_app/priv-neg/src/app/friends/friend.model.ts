
export class Friend {
  id: string;
  picture: string;
  name: string;
  // tieStrength: number;
  // group: string;

  public static FromFBFriend(fbFriend: FBFriend): Friend {
    const f = new Friend();

    f.id = fbFriend.id;
    f.picture = fbFriend.picture.data.url;
    f.name = fbFriend.name;

    return f;
  }

  public static UpdateFromAPIFriend(f: Friend, apiFriend: APIFriend): Friend {
    f.id = apiFriend.to;
    // f.tieStrength = apiFriend.tieStrength;
    // f.group = apiFriend.group;

    return f;
  }
}

export class APIFriend {
  from: string;
  to: string;
  // tieStrength: number;
  // group: string;
}

export class FBFriend {
  id: string;
  name: string;
  picture;
}

