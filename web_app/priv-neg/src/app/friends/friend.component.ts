import {Component, OnInit} from '@angular/core';
import {ActivatedRoute, Router} from '@angular/router';
import 'rxjs/add/operator/switchMap';
import {User} from '../domain/user.model';
import {FriendService} from './friend.service';


@Component({
  selector: 'app-friend',
  templateUrl: './friend.component.html',
})
export class FriendComponent implements OnInit {

  public friend: User;

  constructor(
    private route: ActivatedRoute,
    private friendService: FriendService,
    private router: Router
  ) {}

  ngOnInit() {
    const id = this.route.snapshot.params['id'];
    this.friendService.getUserById(id)
      .then(friend => {
        this.friend = friend;
        console.log(this.friend);
      })
      .catch(this.router.navigate['start'])
    ;
  }
}
