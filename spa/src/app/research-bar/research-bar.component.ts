import {Component, EventEmitter, OnInit, Output} from '@angular/core';
import {AirportApiService} from "../services/airport-api.service";

interface Airport {
  code: string;
}

interface Sensor {
  name: string;
}

interface ApiAirport {
  airport: string;
}

@Component({
  selector: 'app-research-bar',
  templateUrl: './research-bar.component.html',
  styleUrls: ['./research-bar.component.scss']
})
export class ResearchBarComponent implements OnInit {
  @Output() selectedAirportChange: EventEmitter<Airport | undefined> = new EventEmitter<Airport | undefined>();
  @Output() selectedSensorsChange: EventEmitter<Sensor[]> = new EventEmitter<Sensor[]>();
  @Output() rangeDatesChange: EventEmitter<Date[]> = new EventEmitter<Date[]>();

  airports: Airport[] = [];
  selectedAirport: Airport | undefined;
  sensors: Sensor[] = [];
  selectedSensors: Sensor[] = [];
  rangeDates: Date[] = [];

  constructor(private airportService: AirportApiService) {}

  ngOnInit() {
    this.airportService.getAirports().subscribe(
      (data: ApiAirport[]) => {
        this.airports = data.map(item => {
          const airport: Airport = { code: item.airport };
          return airport;
        });
      },
      (error) => {
        console.error('Error fetching airports:', error);
      }
    );
  }

  onAirportSelected() {
    if (this.selectedAirport) {
      this.airportService.getSensorsByAirportIATA(this.selectedAirport.code).subscribe(
        originalSensors => {
          this.sensors = originalSensors.map(originalSensor => ({ name: originalSensor.measurement }));
        },
        error => {
          console.error('Error fetching sensors:', error);
        }
      );
    } else {
      this.sensors = [];
      this.selectedSensorsChange.emit(this.sensors);
    }
    this.selectedAirportChange.emit(this.selectedAirport);
  }

  onSensorsChange() {
    this.selectedSensorsChange.emit(this.selectedSensors);
  }

  onDateRangeChange() {
    this.rangeDatesChange.emit(this.rangeDates);
  }
}
