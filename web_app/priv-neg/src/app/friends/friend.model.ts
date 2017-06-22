
export class Friend {
  id: string;
  picture: string;
  name: string;
  // tieStrength: number;
  clique: string;

  public static FromFBFriend(fbFriend: FBFriend): Friend {
    const f = new Friend();

    f.id = fbFriend.id;
    f.picture = fbFriend.picture.data.url;
    f.name = fbFriend.name;

    return f;
  }

  public static UpdateFromAPIFriend(f: Friend, apiFriend: APIFriend): Friend {
    f.id = apiFriend.id;
    // f.tieStrength = apiFriend.tieStrength;
    // f.clique = apiFriend.clique;

    return f;
  }
}

export class APIFriend {
  id: string;
  clique: string;
  // tieStrength: number;
}

export class FBFriend {
  id: string;
  name: string;
  picture;
}

