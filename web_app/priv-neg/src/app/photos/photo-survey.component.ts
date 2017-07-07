import {Component, OnInit} from '@angular/core';
import {ActivatedRoute, Router} from '@angular/router';
import 'rxjs/add/operator/switchMap';
import {Photo} from '../domain/photo.model';
import {PhotoService} from './photo.service';

@Component({
  selector: 'app-photo-survey',
  templateUrl: './photo-survey.component.html',
})
export class PhotoSurveyComponent implements OnInit {

  public photo: Photo;
  public survey: Survey;

  constructor(
    private route: ActivatedRoute,
    private photoService: PhotoService,
    private router: Router
  ) {}

  ngOnInit() {
    const id = this.route.snapshot.params['id'];
    if (id) {

      this.photo = this.photoService.getPhotoById(id);

      this.survey = new Survey();
      this.survey.photoID = this.photo.id;

      console.log(this.photo);
    }
  }

  submit() {

  }
}

export class Survey {
  photoID: string;

  agreeRecommendation = '';
  agreeRecommendationWhy = '';
  agreeRecommendationChoices = [
    new Option('yes', 'Yes'),
    new Option('no', 'No'),
  ];

  wereYouAskedPermission = [
    new Option('no_doesnt_bother', 'No, but it doesn\'t bother me'),
    new Option('no_and_bother', 'No, and it does bother me'),
    new Option('yes', 'Yes'),
  ];

  tieStrengthDepicted = '';
  tieStrengthDepictedChoices = [
    new Option('yes', 'Yes'),
    new Option('no', 'No'),
  ];
  tieStrengthDepictedWhy: string;

  doYouUnderstandTheMotivesOfOpposition: string;
  doYouUnderstandTheMotivesOfOppositionChoices = [
    new Option('hard-no', 'No, it won\'t change my preference'),
    new Option('soft-no', 'No, but they might change my preference'),
    new Option('yes-no-concede', 'Yes, but I would not concede'),
    new Option('yes-concede-compensation', 'Yes, and I would concede for compensation'),
    new Option('yes-concede', 'Yes, and I would concede'),
  ];

  howWouldYouResolveThisConflict = '';
  howWouldYouResolveThisConflictChoices = [
    new Option('untag', 'I would un-tag people'),
    new Option('remove', 'I would remove the photo entirely'),
    new Option('blur', 'I would blur people from the photo'),
    new Option('crop', 'I would crop people from the photo'),
    new Option('nothing', 'I would do nothing'),
    new Option('other', 'Other')
  ];

}

export class GeneralSurvey {

  // At least every greyscale photo presented by this tool
  // is publicly searchable and viewable by everyone. Has this tool
  // given you more awareness of cases that undesirably breach your privacy
  // i.e. A photo is publicly viewable that you weren't aware of?
  q1Answer = '';
  q1Choices = [
    new Option('yes-no-idea', 'Yes, I had no idea.'),
    new Option('yes-some-idea', 'Yes, I had some awareness but didn\'t realise the extent.'),
    new Option('no', 'No, I was already aware.')
  ]

  // What is your current sharing preference on Facebook
  q2Answer = '';
  q2Choices = [
    new Option('public', 'Everyone'),
    new Option('friends-only', 'Friends Only'),
    new Option('friends-of-friends', 'Friends of friends'),
    new Option('only-me', 'Only me'),
    new Option('dont-know', 'Don\'t know')
  ]

  // How do you usually resolve conflicts with photos that you upload or are tagged in?
  q3Answer = '';
  q3Choices = [
    new Option('untag', 'I would untag myself/others'),
    new Option('crop', 'I would crop the photo'),
    new Option('remove', 'I would remove the photo entirely'),
    new Option('blur', 'I would blur the photo'),
    new Option('nothing', 'I would do nothing'),
    new Option('other', 'Other?')
  ]
  q3Other = '';

  // Do you think that your Facebook profile represents yourself and your relationships well?
  q4Answer = '';
  q4Choices = [
    new Option('yes', 'Yes'),
    new Option('no', 'No')
  ];
  q4Why = '';

  // How much weight would you place on these for quantifying your relationship strength with someone (how close you are to someone)
  q5Parts = [
    new NumberOption('politics', 'Politics'),
    new NumberOption('religion', 'Religion'),
    new NumberOption('work', 'Work'),
    new NumberOption('sports', 'Sports'),
    new NumberOption('family', 'Family member'),
    new NumberOption('location', 'Location'),
    new NumberOption('education', 'Education'),
    new NumberOption('favourite-teams', 'Favourite Teams'),
    new NumberOption('inspirational-people', 'Inspirational People'),
    new NumberOption('languages', 'Languages'),
    new NumberOption('music', 'Music'),
    new NumberOption('movies', 'Movies'),
    new NumberOption('likes', 'Likes'),
    new NumberOption('groups', 'Groups'),
    new NumberOption('events', 'Events'),
  ];

  // How could this tool be improved?
  q6Answer = ''

  // Any further comments?
  q7Answer = '';
}

export class PhotoSurvey {

}

export class Option {
  id: string;
  description: string;

  constructor(id: string, description: string) {
    this.id = id;
    this.description = description;
  }
}

export class NumberOption {
  id: string;
  description: string;
  value: number;

  constructor(id: string, description: string) {
    this.id = id;
    this.description = description;
  }
}
