import {Injectable} from '@angular/core';
import {FriendService} from '../friends/friend.service';
import {APIConflict, APIPhoto, APIReason, Conflict, FbGraphPhoto, Photo, Reason} from '../domain/photo.model';

@Injectable()
export class PhotoResolver {

  constructor(
    private friendService: FriendService
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
    photo.pending = apiPhoto.pending;
    photo.categories = apiPhoto.categories;

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

}
