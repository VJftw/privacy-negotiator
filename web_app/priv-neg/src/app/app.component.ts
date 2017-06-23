import {Component, OnInit} from '@angular/core';
import { environment } from '../environments/environment';
import {AuthService, SessionUser} from './auth.service';
import {APIService} from './api.service';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styles: [`
    .info a div {
      font-size: 10px;
      line-height: 32px;
    }
  `],
})
export class AppComponent implements OnInit  {
  version = environment.version;
  latency = '';

  constructor(
    public authService: AuthService,
    private apiService: APIService
  ) {}

  ngOnInit() {
    this.updateLatency();
  }

  private updateLatency() {
    const timeStart = performance.now();

    this.apiService.get('/v1/health').then(res => {
      const timeEnd = performance.now();
      this.latency = ('   ' + (timeEnd - timeStart).toFixed(0)).slice(-3);
    });

    this.sleep(5000).then(() => this.updateLatency());
  }

  private sleep(time) {
    return new Promise((resolve) => setTimeout(resolve, time));
  }
}
