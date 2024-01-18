import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';

import { AppComponent } from './app.component';
import { ResearchBarComponent } from './research-bar/research-bar.component';
import {ThemeService} from "../theme/service";

@NgModule({
  declarations: [
    AppComponent,
    ResearchBarComponent
  ],
  imports: [
    BrowserModule
  ],
  providers: [ThemeService],
  bootstrap: [AppComponent]
})
export class AppModule { }
