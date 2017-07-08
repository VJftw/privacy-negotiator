import {Component, OnInit} from '@angular/core';
import {ActivatedRoute} from '@angular/router';
import 'rxjs/add/operator/switchMap';
import {Photo} from '../domain/photo.model';
import {PhotoService} from './photo.service';
import {APISurvey, Option} from './survey.component';
import {APIService} from '../api.service';

@Component({
  selector: 'app-photo-survey',
  templateUrl: './photo-survey.component.html',
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
export class PhotoSurveyComponent implements OnInit {

  public photo: Photo;
  public survey: PhotoSurvey;
  public submitted = false;

  constructor(
    private route: ActivatedRoute,
    private photoService: PhotoService,
    private apiService: APIService,
  ) {}

  ngOnInit() {
    const id = this.route.snapshot.params['id'];
    if (id) {

      this.photoService.getPhotoById(id).then(photo => {
        this.photo = photo;
        this.survey = new PhotoSurvey(this.photo);
        console.log(this.photo);
      });

    }
  }

  submit() {
    this.submitted = true;
    const apiSurvey = new APISurvey();
    apiSurvey.type = 'photo';
    apiSurvey.photoID = this.photo.id;
    apiSurvey.data = this.survey;

    this.apiService.post('/v1/surveys', apiSurvey).then(res => console.log(res));
  }
}


export class PhotoSurvey {
  photoId: string;

  // Were you asked for your permission before the photo was uploaded?
  q1Answer = '';
  q1Choices = [
    new Option('yes', 'Yes'),
    new Option('no-dont-mind', 'No, but I don\'t mind'),
    new Option('no-do-mind', 'No, and I do mind')
  ];

  // Do you agree with the recommendation?
  q2Answer = '';
  q2Choices = [
    new Option('yes', 'Yes'),
    new Option('no', 'No')
  ];

  // Do you understand the motives of the opposing parties in the conflict?
  q3Answer = '';
  q3Choices = [
    new Option('yes-no-concede', 'Yes, but I would not concede my preference'),
    new Option('yes-concede', 'Yes, and I would concede my preference'),
    new Option('no-may-concede', 'No, but I may concede knowing their reasons'),
    new Option('no-no-concede', 'No, and I would not concede')
  ];

  constructor(photo: Photo) {
    this.photoId = photo.id;
  }
}
