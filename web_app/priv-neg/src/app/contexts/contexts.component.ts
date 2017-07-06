import {AfterViewChecked, Component, OnInit} from '@angular/core';
import { ContextService } from './context.service';
import {Context} from '../domain/context.model';

declare var Materialize: any;

@Component({
  selector: 'app-contexts',
  templateUrl: './contexts.component.html'
})
export class ContextsComponent implements OnInit, AfterViewChecked {

  newContext: Context;
  editing: boolean;

  constructor(
    public contextService: ContextService,
  ) {
    this.newContext = new Context('', true);
    this.editing = false;
  }

  ngOnInit() {
    this.contextService.updateContextsFromAPI();
  }

  getGlobalContexts(): Context[] {
    const cats: Context[] = [];
    for (const cat of this.contextService.getContexts()) {
      if (!cat.personal) {
        cats.push(cat);
      }
    }

    return cats;
  }

  getPersonalContexts(): Context[] {
    const cats: Context[] = [];
    for (const cat of this.contextService.getContexts()) {
      if (cat.personal) {
        cats.push(cat);
      }
    }

    return cats;
  }

  toggleEdit() {
    if (!this.editing) {
      this.editing = true;
    } else {
      this.editing = false;
      if (this.newContext.name.length > 0) {
        this.contextService.createContext(this.newContext.name).catch(
          () => console.log('Invalid Context')
        ).then(
          () => this.newContext = new Context('', true)
        );
      }
    }
  }

  ngAfterViewChecked() {
    Materialize.updateTextFields();
  }

}
