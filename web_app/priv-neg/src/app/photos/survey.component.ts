import {Component, OnInit, OnDestroy} from '@angular/core';
import {NavigationEnd, Router} from '@angular/router';

import 'rxjs/add/operator/switchMap';
import {Photo} from '../domain/photo.model';
import {PhotoService} from './photo.service';
import {APIService} from '../api.service';

@Component({
  selector: 'app-survey',
  templateUrl: './survey.component.html',
  styles: [`
[type="radio"]:checked + label:after, [type="radio"].with-gap:checked + label:before, [type="radio"].with-gap:checked + label:after {
  background-color: #2196f3 !important;
}

[type="radio"]:checked + label:after, [type="radio"].with-gap:checked + label:before, [type="radio"].with-gap:checked + label:after {
  border: 2px solid #2196f3 !important;
}

[type="checkbox"]:checked + label:before {
  border-right: 2px solid #2196f3;
  border-bottom: 2px solid #2196f3;
}

textarea.materialize-textarea:focus:not([readonly]) {
  border-bottom: 1px solid #2196f3;
  box-shadow: 0 1px 0 0 #2196f3;
}

textarea.materialize-textarea:focus:not([readonly]) + label {
  color: #2196f3;
}
  `]
})
export class SurveyComponent {

  public survey = new GeneralSurvey();
  public submitted = false;

  constructor(
    private router: Router,
    private apiService: APIService,
  ) {}

  private extendedValidation(): boolean {
    this.survey.q3Answer = [];
    for (const q3Checkbox of this.survey.q3Checkboxes) {
      if (q3Checkbox.selected) {
        this.survey.q3Answer.push(q3Checkbox.description)
      }
    }

    if (this.survey.q3Answer.length <= 0) {
      return false;
    }

    if (this.survey.q3Checkboxes[5].selected && this.survey.q3Other.trim().length <= 5) {
      return false;
    }

    if (this.survey.q4Answer === 'No' && this.survey.q4Why.trim().length <= 5) {
      return false;
    }

    if (this.survey.q5Answer === 'No' && this.survey.q5Why.trim().length <= 5) {
      return false;
    }

    return true;

  }

  submit() {
    this.submitted = true;
    const apiSurvey = new APISurvey();
    apiSurvey.type = 'general';
    apiSurvey.data = this.survey;

    this.apiService.post('/v1/surveys', apiSurvey).then(res => console.log(res));
  }
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
  ];

  // What is your current sharing preference on Facebook
  q2Answer = '';
  q2Choices = [
    new Option('public', 'Everyone'),
    new Option('friends-only', 'Friends Only'),
    new Option('only-me', 'Only me'),
    new Option('dont-know', 'I Don\'t know')
  ];

  // How do you usually resolve conflicts with photos that you upload or are tagged in?
  q3Answer = [];
  q3Checkboxes = [
    new CheckboxOption('untag', 'I would untag myself/others'),
    new CheckboxOption('crop', 'I would crop the photo'),
    new CheckboxOption('remove', 'I would remove the photo entirely'),
    new CheckboxOption('blur', 'I would blur the photo'),
    new CheckboxOption('nothing', 'I would do nothing'),
    new CheckboxOption('other', 'Other?'),
  ];
  q3Other = '';

  // Do you think that your Facebook profile represents yourself and your relationships well?
  q4Answer = '';
  q4Choices = [
    new Option('yes', 'Yes'),
    new Option('no', 'No')
  ];
  q4Why = '';

  q5Answer = '';
  q5Choices = [
    new Option('yes', 'Yes'),
    new Option('no', 'No')
  ];
  q5Why = '';

  // How much weight would you place on these for quantifying your relationship strength with someone (how close you are to someone)
  q6Parts = [
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
  q7Answer = '';

  // Any further comments?
  q8Answer = '';

}

export class Option {
  id: string;
  description: string;

  constructor(id: string, description: string) {
    this.id = id;
    this.description = description;
  }
}

export class CheckboxOption {
  id: string;
  description: string;
  selected = false;

  constructor(id: string, description: string) {
    this.id = id;
    this.description = description;
  }
}

export class NumberOption {
  id: string;
  description: string;
  value;

  constructor(id: string, description: string) {
    this.id = id;
    this.description = description;
  }
}

export class APISurvey {
  type: string;
  photoID: string;
  data;
}
