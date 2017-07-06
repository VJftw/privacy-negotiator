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
import { CategoryService } from './categories/category.service';
import { FriendService } from './friends/friend.service';
import {PhotoResolver} from './photos/photo.resolver';
import {CategoriesComponent} from './categories/categories.component';
import {SessionService} from './session.service';
import {SurveyComponent} from "./photos/survey.component";

@NgModule({
  declarations: [
    AppComponent,
    IndexComponent,
    CategoriesComponent,
    PhotosComponent,
    PhotoDetailComponent,
    FriendsComponent,
    SurveyComponent,
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
      { path: 'categories', component: CategoriesComponent, canActivate: [SessionService] },
      { path: 'photos/:id', component: PhotoDetailComponent, canActivate: [SessionService] },
      { path: 'survey/:id', component: SurveyComponent, canActivate: [SessionService] },
      { path: 'friends', component: FriendsComponent, canActivate: [SessionService] }
    ], { useHash: true })
  ],
  providers: [
    APIService,
    WebSocketService,
    AuthService,
    SessionService,
    FacebookService,
    PhotoService,
    CategoryService,
    FriendService,
    PhotoResolver,
  ],
  bootstrap: [AppComponent]
})
export class AppModule { }
