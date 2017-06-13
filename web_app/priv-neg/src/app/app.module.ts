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
import { CommunitiesComponent } from './communities/communities.component';

import { FacebookService } from 'ngx-facebook';
import { APIService } from './api.service';
import { WebSocketService } from './websocket.service';
import { AuthService } from './auth.service';
import { PhotoService } from './photos/photo.service';

@NgModule({
  declarations: [
    AppComponent,
    IndexComponent,
    PhotosComponent,
    PhotoDetailComponent,
    CommunitiesComponent
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
      { path: 'photos', component: PhotosComponent, canActivate: [AuthService] },
      { path: 'photos/:id', component: PhotoDetailComponent, canActivate: [AuthService] },
      { path: 'communities', component: CommunitiesComponent, canActivate: [AuthService] }
    ], { useHash: true })
  ],
  providers: [
    APIService,
    WebSocketService,
    AuthService,
    FacebookService,
    PhotoService
  ],
  bootstrap: [AppComponent]
})
export class AppModule { }
