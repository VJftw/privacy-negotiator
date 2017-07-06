import {Component, OnInit} from '@angular/core';
import {ActivatedRoute, Router} from '@angular/router';
import 'rxjs/add/operator/switchMap';
import {PhotoService} from './photo.service';
import {Photo} from '../domain/photo.model';
import {Context, ContextSelection} from '../domain/context.model';
import {ContextService} from '../contexts/context.service';
import {User} from '../domain/user.model';


@Component({
  selector: 'app-photo',
  templateUrl: './photo-detail.component.html',
})
export class PhotoDetailComponent implements OnInit {

  public photo: Photo;
  public editing = false;
  public showConflictHelp = false;
  private contextSelection: Map<string, ContextSelection> = new Map();

  constructor(
    private route: ActivatedRoute,
    private photoService: PhotoService,
    private contextService: ContextService,
    private router: Router
  ) {}

  ngOnInit() {
    const id = this.route.snapshot.params['id'];
    this.photo = this.photoService.getPhotoById(id);

    if (!this.photo) {
      this.router.navigate(['start']);
    }

    console.log(this.photo);
    this.updateChoices();
  }

  public getContexts(): ContextSelection[] {
    return Array.from(this.contextSelection.values());
  }

  public updateChoices() {
    for (const cat of this.contextService.getContexts()) {
      if (Context.isIn(cat, this.photo.contexts)) {
        this.contextSelection.set(cat.name, new ContextSelection(cat.name, true));
      } else {
        this.contextSelection.set(cat.name, new ContextSelection(cat.name, false));
      }
    }
  }

  toggleEdit() {
    this.editing = !this.editing;
    if (!this.editing) {
      this.photo.contexts = [];
      for (const contextSelection of this.getContexts()) {
        if (contextSelection.isActive) {
          this.photo.contexts.push(this.contextService.getContext(contextSelection.name));
        }
      }
      console.log(this.photo);
      this.photoService.updatePhoto(this.photo);
    }
  }

  isInConflict(taggedUser: User): boolean {
    for (const conflict of this.photo.conflicts) {
      for (const partyUser of conflict.parties) {
        if (partyUser.id === taggedUser.id) {
          return true;
        }
      }
    }

    return false;
  }

  toggleConflictHelp() {
    this.showConflictHelp = !this.showConflictHelp;
  }


}


