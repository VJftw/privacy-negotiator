import {Injectable} from '@angular/core';
import {APIService} from '../api.service';
import {APIContext, Context} from '../domain/context.model';

@Injectable()
export class ContextService {

  protected contexts: Map<string, Context>;

  constructor(
    private apiService: APIService
  ) {
    this.contexts = new Map();
  }

  public getContexts(): Context[] {
    return Array.from(this.contexts.values());
  }

  public getContext(name: string) {
    return this.contexts.get(name);
  }

  public updateContextsFromAPI(): Promise<void> {
    const uri = '/v1/categories';

    return this.apiService.get(uri).then(response => {
      for (const cat of response.json() as APIContext[]) {
        this.contexts.set(cat.name, Context.fromAPIContext(cat));
      }
    });
  }

  public createContext(c: string): Promise<any> {
    const uri = '/v1/categories';

    return this.apiService.post(uri, new Context(c)).then(response => {
      console.log(response);
      this.contexts.set(c, new Context(c, true));
    });
  }
}
