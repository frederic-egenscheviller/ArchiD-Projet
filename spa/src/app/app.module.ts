import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';
import { DropdownModule } from 'primeng/dropdown';

import { AppComponent } from './app.component';
import { ResearchBarComponent } from './research-bar/research-bar.component';
import {ThemeService} from "../theme/service";
import {FormsModule} from "@angular/forms";
import {BrowserAnimationsModule} from "@angular/platform-browser/animations";
import {MultiSelectModule} from "primeng/multiselect";
import {AirportSelectorComponent} from "./research-bar/airport-selector/airport-selector.component";
import {SensorSelectorComponent} from "./research-bar/sensor-selector/sensor-selector.component";

@NgModule({
  declarations: [
    AppComponent,
    ResearchBarComponent
  ],
  imports: [
    BrowserModule,
    DropdownModule,
    MultiSelectModule,
    BrowserAnimationsModule,
    FormsModule,
    AirportSelectorComponent,
    SensorSelectorComponent
  ],
  providers: [ThemeService],
  bootstrap: [AppComponent]
})
export class AppModule { }
