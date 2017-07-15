import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { RouterModule } from '@angular/router';
import { HttpModule } from '@angular/http';

import { MaterializeModule } from 'angular2-materialize';
import { FacebookModule } from 'ngx-facebook';
import { InfiniteScrollModule } from 'ngx-infinite-scroll';

import { AppComponent } from './app.component';
import { IndexComponent } from './index/index.component';
import { PhotosComponent } from './photos/photos.component';
import { PhotoDetailComponent } from './photos/photo-detail.component';
import { FriendsComponent } from './friends/friends.component';

import { FacebookService } from 'ngx-facebook';
import { APIService } from './api.service';
import { WebSocketService } from './websocket.service';
import { AuthService } from './auth.service';
import { PhotoService } from './photos/photo.service';
import { ContextService } from './contexts/context.service';
import { FriendService } from './friends/friend.service';
import {PhotoResolver} from './photos/photo.resolver';
import {ContextsComponent} from './contexts/contexts.component';
import {SessionService} from './session.service';
import {SurveyComponent} from './photos/survey.component';
import {PhotoSurveyComponent} from './photos/photo-survey.component';
import {PrivacyPolicyComponent} from './index/privacy-policy.component';
import {FriendComponent} from './friends/friend.component';

@NgModule({
  declarations: [
    AppComponent,
    IndexComponent,
    ContextsComponent,
    PhotosComponent,
    PhotoDetailComponent,
    FriendsComponent,
    FriendComponent,
    SurveyComponent,
    PhotoSurveyComponent,
    PrivacyPolicyComponent,
  ],
  imports: [
    BrowserModule,
    InfiniteScrollModule,
    FormsModule,
    HttpModule,
    MaterializeModule,
    FacebookModule.forRoot(),
    RouterModule.forRoot([
      { path: '', redirectTo: '/start', pathMatch: 'full' },
      { path: 'start', component: IndexComponent },
      { path: 'photos', component: PhotosComponent, canActivate: [SessionService] },
      { path: 'contexts', component: ContextsComponent, canActivate: [SessionService] },
      { path: 'survey/:id', component: PhotoSurveyComponent, canActivate: [SessionService] },
      { path: 'survey', component: SurveyComponent, canActivate: [SessionService] },
      { path: 'photos/:id', component: PhotoDetailComponent, canActivate: [SessionService] },
      { path: 'friends/:id', component: FriendComponent, canActivate: [SessionService] },
      { path: 'friends', component: FriendsComponent, canActivate: [SessionService] },
      { path: 'privacy-policy', component: PrivacyPolicyComponent}
    ], { useHash: true })
  ],
  providers: [
    APIService,
    WebSocketService,
    AuthService,
    SessionService,
    FacebookService,
    PhotoService,
    ContextService,
    FriendService,
    PhotoResolver,
  ],
  bootstrap: [AppComponent]
})
export class AppModule { }
