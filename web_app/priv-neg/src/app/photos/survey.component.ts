import {Component, OnInit} from '@angular/core';
import {ActivatedRoute, Router} from '@angular/router';
import 'rxjs/add/operator/switchMap';
import {Photo} from '../domain/photo.model';
import {PhotoService} from "./photo.service";

@Component({
  selector: 'app-survey',
  templateUrl: './survey.component.html',
})
export class SurveyComponent implements OnInit {

  public photo: Photo;
  public survey: Survey;

  constructor(
    private route: ActivatedRoute,
    private photoService: PhotoService,
    private router: Router
  ) {}

  ngOnInit() {
    const id = this.route.snapshot.params['id'];
    this.photo = this.photoService.getPhotoById(id);

    if (!this.photo) {
      this.router.navigate(['start']);
    }

    this.survey = new Survey();
    this.survey.photoID = this.photo.id;

    console.log(this.photo);
  }

  submit() {

  }
}

export class Survey {
  photoID: string;

  howWouldYouResolveThisConflict = '';
  howWouldYouResolveThisConflictChoices = [
    new Option('untag', 'I would un-tag people'),
    new Option('remove', 'I would remove the photo entirely'),
    new Option('blur', 'I would blur people from the photo'),
    new Option('crop', 'I would crop people from the photo'),
    new Option('nothing', 'I would do nothing'),
    new Option('other', 'Other')
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


  howDoYouUsuallyResolveConflicts = '';

  wereYouAskedPermission = [
    new Option('no_doesnt_bother', 'No, but it doesn\'t bother me'),
    new Option('no_and_bother', 'No, and it does bother me'),
    new Option('yes', 'Yes'),
  ];

}

export class Option {
  id: string;
  description: string;

  constructor(id: string, description: string) {
    this.id = id;
    this.description = description;
  }
}
