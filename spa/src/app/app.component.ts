import {Component, OnInit} from '@angular/core';
import {ThemeService} from "./services/theme";
import {AirportApiService} from "./services/airport-api.service";
import { DatePipe } from '@angular/common';

interface Airport {
  code: string;
}

interface Sensor {
  name: string;
}

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent implements OnInit {
  title = 'Airport MQTT Dashboard';

  data: any;

  options: any;
  constructor(private themeService: ThemeService, private airportService: AirportApiService, private datePipe: DatePipe) {}

  selectedAirport: Airport | undefined;
  selectedSensors: Sensor[] = [];
  rangeDates: Date[] = [];

  onSelectedAirportChange(selectedAirport: Airport | undefined) {
    this.selectedAirport = selectedAirport;
    this.updateData();
  }

  onSelectedSensorsChange(selectedSensors: Sensor[]) {
    this.selectedSensors = selectedSensors;
    this.updateData();
  }

  onRangeDatesChange(rangeDates: Date[]) {
    this.rangeDates = rangeDates;
    this.updateData();
  }

  ngOnInit(): void {
    this.themeService.toggleBodyClass(this.themeService.isDarkMode());
    this.themeService.watchDarkMode((darkMode: boolean) => {
      this.themeService.toggleBodyClass(darkMode);
    });
  }

  updateData() {
    this.data = [];
    this.options = [];

    const documentStyle = getComputedStyle(document.documentElement);
    const textColor = documentStyle.getPropertyValue('--text-color');
    const textColorSecondary = documentStyle.getPropertyValue('--text-color-secondary');
    const surfaceBorder = documentStyle.getPropertyValue('--surface-border');

    let startDate = this.datePipe.transform(this.rangeDates[0], 'yyyy-MM-ddTHH:mm:ssZ');
    if (startDate != null && this.selectedSensors != null) {
      let endDate = this.datePipe.transform(this.rangeDates[this.rangeDates.length - 1], 'yyyy-MM-ddTHH:mm:ssZ');
      let endDatedate;
      if (endDate === null) {
        endDatedate = new Date(startDate!);
        endDatedate!.setDate(endDatedate!.getDate() + 1);
        endDate = this.datePipe.transform(endDatedate, 'yyyy-MM-ddTHH:mm:ssZ');
      }
      endDate = endDate?.split('+')[0] + 'Z';
      startDate = startDate?.split('+')[0] + 'Z';

      const datasets: { label: string; data: number[]; fill: boolean; borderColor: string; tension: number; }[] = [];

      this.selectedSensors.forEach(sensor => {
        this.airportService.getMeasurementDataByDateRangeAndType(
          this.selectedAirport!.code,
          startDate!,
          endDate!,
          sensor.name
        ).subscribe(
          measurementData => {
            datasets.push({
              label: sensor.name,
              data: measurementData.map(data => data.value),
              fill: false,
              borderColor: this.getRandomColor(),
              tension: 0.4,
            });

            if (datasets.length === this.selectedSensors.length) {
              this.data = {
                labels: measurementData.map(data => data.time),
                datasets: datasets,
              };
            }
          },
          error => {
            console.error('Error fetching measurement data:', error);
          }
        );
      });
    }
    this.options = {
      maintainAspectRatio: false,
      aspectRatio: 0.6,
      plugins: {
        legend: {
          labels: {
            color: textColor
          }
        }
      },
      scales: {
        x: {
          ticks: {
            color: textColorSecondary
          },
          grid: {
            color: surfaceBorder,
            drawBorder: false
          }
        },
        y: {
          ticks: {
            color: textColorSecondary
          },
          grid: {
            color: surfaceBorder,
            drawBorder: false
          }
        }
      }
    };
  }

  getRandomColor() {
    const letters = '0123456789ABCDEF';
    let color = '#';
    for (let i = 0; i < 6; i++) {
      color += letters[Math.floor(Math.random() * 16)];
    }
    return color;
  }
}
