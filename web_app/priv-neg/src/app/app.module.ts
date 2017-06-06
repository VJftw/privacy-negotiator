import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { RouterModule } from '@angular/router';
import { HttpModule } from '@angular/http';

import { MaterializeModule } from 'angular2-materialize';
import { FacebookModule } from 'ngx-facebook';

import { AppComponent } from './app.component';
import { IndexComponent } from './index/index.component';
import { PhotosComponent } from './photos/photos.component';

@NgModule({
  declarations: [
    AppComponent,
    IndexComponent,
    PhotosComponent
  ],
  imports: [
    BrowserModule,
    FormsModule,
    HttpModule,
    MaterializeModule,
    FacebookModule.forRoot(),
    RouterModule.forRoot([
      { path: '', redirectTo: '/start', pathMatch: 'full' },
      { path: 'start', component: IndexComponent },
      { path: 'photos', component: PhotosComponent }
    ])
  ],
  providers: [],
  bootstrap: [AppComponent]
})
export class AppModule { }
