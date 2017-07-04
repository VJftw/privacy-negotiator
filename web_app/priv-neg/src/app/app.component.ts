import {Component, OnInit} from '@angular/core';
import { environment } from '../environments/environment';
import {AuthService} from './auth.service';
import {APIService} from './api.service';
import {SessionService} from './session.service';

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
  queueSize: number;
  apiStatus = 'OFFLINE';
  info = 0;

  constructor(
    public sessionService: SessionService,
    private apiService: APIService
  ) {}

  ngOnInit() {
    this.updateLatency();
  }

  public cycleInfo() {
    if (this.info < 2) {
      this.info++;
    } else {
      this.info = 0;
    }
  }

  private updateLatency() {
    const timeStart = performance.now();

    this.apiService.get('/v1/health').then(res => {
      const apiHealth = res.json() as ApiHealth;
      const timeEnd = performance.now();
      this.latency = ('   ' + (timeEnd - timeStart).toFixed(0)).slice(-3);
      this.queueSize = apiHealth.queueSize;

      if ((timeEnd - timeStart) < 500 && apiHealth.queueSize < 20) {
        this.apiStatus = 'OK';
      } else {
        this.apiStatus = 'BUSY';
      }
    }).catch(() => {
      this.apiStatus = 'OFFLINE';
    });

    this.sleep(5000).then(() => this.updateLatency());
  }

  private sleep(time) {
    return new Promise((resolve) => setTimeout(resolve, time));
  }
}

class ApiHealth {
  queueSize: number;
}
