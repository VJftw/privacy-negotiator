import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { RouterModule }   from '@angular/router';
import { HttpModule } from '@angular/http';

import { MaterializeModule } from 'angular2-materialize';
import { FacebookModule } from 'ngx-facebook';

import { AppComponent } from './app.component';
import {IndexComponent} from './index/index.component';

@NgModule({
  declarations: [
    AppComponent,
    IndexComponent
  ],
  imports: [
    BrowserModule,
    FormsModule,
    HttpModule,
    MaterializeModule,
    FacebookModule.forRoot(),
    RouterModule.forRoot([
      {
        path: '',
        component: IndexComponent
      }
    ])
  ],
  providers: [],
  bootstrap: [AppComponent]
})
export class AppModule { }
