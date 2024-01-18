import {Component, OnInit} from '@angular/core';
import {ThemeService} from "../theme/service";

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent implements OnInit {
  title = 'Airport MQTT Dashboard';
  constructor(private themeService: ThemeService) {}

  ngOnInit(): void {
    this.themeService.toggleBodyClass(this.themeService.isDarkMode());
    this.themeService.watchDarkMode((darkMode: boolean) => {
      this.themeService.toggleBodyClass(darkMode);
    });
  }
}
