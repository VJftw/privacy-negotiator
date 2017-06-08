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

import { FacebookService } from 'ngx-facebook';
import { AuthService } from './auth.service';

@NgModule({
  declarations: [
    AppComponent,
    IndexComponent,
    PhotosComponent
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
      { path: 'photos', component: PhotosComponent, canActivate: [AuthService] }
    ], { useHash: true })
  ],
  providers: [
    AuthService,
    FacebookService
  ],
  bootstrap: [AppComponent]
})
export class AppModule { }
