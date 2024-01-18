import {Component, Input, OnInit} from '@angular/core';
import {MultiSelectModule} from "primeng/multiselect";
import {FormsModule} from "@angular/forms";

interface Sensor {
  name: string;
}

@Component({
  selector: 'app-sensor-selector',
  standalone: true,
  imports: [
    MultiSelectModule,
    FormsModule
  ],
  templateUrl: './sensor-selector.component.html',
  styleUrl: './sensor-selector.component.scss'
})
export class SensorSelectorComponent implements OnInit{
  @Input()
  airportIsSelected = false;

  sensors: Sensor[] | undefined;

  selectedSensors: Sensor[] | undefined;

  ngOnInit() {
    this.sensors = [
      { name: 'temperature' },
      { name: 'wind' },
      { name: 'pressure' }
    ];
  }
}
