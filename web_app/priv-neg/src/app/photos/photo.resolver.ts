import {Injectable} from '@angular/core';
import {FriendService} from '../friends/friend.service';
import {APIConflict, APIPhoto, APIReason, Conflict, FbGraphPhoto, Photo, Reason} from '../domain/photo.model';
import {ContextService} from '../contexts/context.service';
import {SessionService} from '../session.service';

@Injectable()
export class PhotoResolver {

  constructor(
    private friendService: FriendService,
    private contextService: ContextService,
    private sessionService: SessionService,
  ) {}

  public photoFromFBPhoto(fbPhoto: FbGraphPhoto): Photo {
    const p = new Photo(fbPhoto.id);

    p.createdTime = fbPhoto.created_time;
    p.url = fbPhoto.images[0].source;
    p.albumId = fbPhoto.album.id;
    p.from = fbPhoto.from;

    return p;
  }

  public photoUpdateFromAPIPhoto(photo: Photo, apiPhoto: APIPhoto): Photo {
    console.log(apiPhoto);
    photo.pending = apiPhoto.pending;
    photo.contexts = [];

    for (const cat of apiPhoto.categories) {
      photo.contexts.push(this.contextService.getContext(cat));
    }

    if (apiPhoto.userCategories[this.sessionService.getUser().id]) {
      for (const cat of apiPhoto.userCategories[this.sessionService.getUser().id]) {
        photo.contexts.push(this.contextService.getContext(cat));
      }
    }

    photo.taggedUsers = [];
    for (const taggedUserId of apiPhoto.taggedUsers) {
      this.friendService.getUserById(taggedUserId).then(
        u => photo.taggedUsers.push(u)
      );
    }

    photo.allowedUsers = [];
    for (const allowedUserId of apiPhoto.allowedUsers) {
      this.friendService.getUserById(allowedUserId).then(
        u => photo.allowedUsers.push(u)
      );
    }

    photo.blockedUsers = [];
    for (const blockedUserId of apiPhoto.blockedUsers) {
      this.friendService.getUserById(blockedUserId).then(
        u => photo.blockedUsers.push(u)
      );
    }

    photo.conflicts = [];
    for (const apiConflict of apiPhoto.conflicts) {
      photo.conflicts.push(this.conflictFromAPIConflict(apiConflict));
    }

    return photo;
  }

  public conflictFromAPIConflict(apiConflict: APIConflict): Conflict {
    const conflict = new Conflict(apiConflict.id);

    this.friendService.getUserById(apiConflict.target).then(
      u => conflict.target = u
    );

    for (const userID of apiConflict.parties) {
      this.friendService.getUserById(userID).then(
        u => conflict.parties.push(u)
      );
    }

    for (const apiReason of apiConflict.reasoning) {
      conflict.reasoning.push(this.reasonFromAPIReason(apiReason));
    }

    conflict.result = apiConflict.result;

    return conflict;
  }

  public reasonFromAPIReason(apiReason: APIReason): Reason {
    const reason = new Reason();

    this.friendService.getUserById(apiReason.id).then(
      u => reason.user = u
    );

    reason.vote = apiReason.vote;

    return reason;
  }

  public APIPhotoFromPhoto(p: Photo): APIPhoto {
    const apiPhoto = new APIPhoto();

    apiPhoto.id = p.id;
    apiPhoto.categories = [];
    apiPhoto.userCategories = new Map();
    apiPhoto.userCategories[this.sessionService.getUser().id] = [];
    for (const cat of p.contexts) {
      if (cat.personal) {
        const currentUserCats = apiPhoto.userCategories[this.sessionService.getUser().id];
        currentUserCats.push(cat.name);
        apiPhoto.userCategories[this.sessionService.getUser().id] = currentUserCats;
      } else {
        apiPhoto.categories.push(cat.name);
      }
    }

    return apiPhoto;

  }

}
