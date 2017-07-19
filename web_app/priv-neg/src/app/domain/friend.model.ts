export class APIFriend {
  id: string;
  cliques: string[];
  tieStrength: number;

  tieStrengthDetails; // vectorName: amount of similarity
}

export class TieStrengthDetail {
  public static VALS = {
    'gender': {
      humanKey: 'Same gender',
      bool: true,
    },
    'ageRange': {
      humanKey: 'Same age group',
      bool: true,
    },
    'family': {
      humanKey: 'Family member',
      bool: true,
    },
    'hometown': {
      humanKey: 'Same hometown',
      bool: true,
    },
    'location': {
      humanKey: 'Same current location',
      bool: true,
    },
    'education': {
      humanKey: 'Educational institutions in common',
      bool: false,
    },
    'events': {
      humanKey: 'Events in common',
      bool: false,
    },
    // 'groups': {
    //   humanKey: 'Groups in common',
    //   bool: false,
    // },
    'favouriteTeams': {
      humanKey: 'Favourite teams in common',
      bool: false,
    },
    'inspirationalPeople': {
      humanKey: 'Inspirational people in common',
      bool: false,
    },
    'languages': {
      humanKey: 'Languages in common',
      bool: false,
    },
    'sports': {
      humanKey: 'Sports in common',
      bool: false,
    },
    'work': {
      humanKey: 'Work places in common',
      bool: false,
    },
    'music': {
      humanKey: 'Music in common',
      bool: false,
    },
    'movies': {
      humanKey: 'Movies in common',
      bool: false,
    },
    'likes': {
      humanKey: 'Likes in common',
      bool: false,
    },
    'political': {
      humanKey: 'Same political alignment',
      bool: true,
    },
    'religion': {
      humanKey: 'Same religion',
      bool: true,
    }
  };

  key: string;
  humanKey: string;
  bool: boolean;
  value: number;

  constructor(key: string, humanKey: string, bool: boolean, value: number) {
    this.key = key;
    this.humanKey = humanKey;
    this.bool = bool;
    this.value = value;
  }
}
