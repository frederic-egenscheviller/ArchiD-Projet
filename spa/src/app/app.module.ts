import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';
import { DropdownModule } from 'primeng/dropdown';
import { AppComponent } from './app.component';
import { ResearchBarComponent } from './research-bar/research-bar.component';
import {ThemeService} from "./services/theme";
import {FormsModule} from "@angular/forms";
import {BrowserAnimationsModule} from "@angular/platform-browser/animations";
import {MultiSelectModule} from "primeng/multiselect";
import {HttpClientModule} from "@angular/common/http";
import {TabViewModule} from "primeng/tabview";
import {ChartModule} from "primeng/chart";
import {CalendarModule} from "primeng/calendar";
import {DatePipe} from "@angular/common";

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
    HttpClientModule,
    TabViewModule,
    ChartModule,
    CalendarModule,
  ],
  providers: [
    ThemeService,
    DatePipe
  ],
  bootstrap: [AppComponent]
})
export class AppModule { }
