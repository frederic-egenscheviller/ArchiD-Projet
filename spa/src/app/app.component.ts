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
    this.updateData(false);
  }

  onSelectedSensorsChange(selectedSensors: Sensor[]) {
    this.selectedSensors = selectedSensors;
    this.updateData(false);
  }

  onRangeDatesChange(rangeDates: Date[]) {
    this.rangeDates = rangeDates;
    this.updateData(false);
  }

  onSelectedAirportChangeAverage(selectedAirport: Airport | undefined) {
    this.selectedAirport = selectedAirport;
    this.updateData(true);
  }

  onSelectedSensorsChangeAverage(selectedSensors: Sensor[]) {
    this.selectedSensors = selectedSensors;
    this.updateData(true);
  }

  onRangeDatesChangeAverage(rangeDates: Date[]) {
    this.rangeDates = rangeDates;
    this.updateData(true);
  }

  ngOnInit(): void {
    this.themeService.toggleBodyClass(this.themeService.isDarkMode());
    this.themeService.watchDarkMode((darkMode: boolean) => {
      this.themeService.toggleBodyClass(darkMode);
    });
  }

  loadPage() {
    this.data = [];
    this.options = []
    this.selectedAirport = undefined;
    this.selectedSensors = [];
    this.rangeDates = [];
  }

  updateData(average: boolean){

    const documentStyle = getComputedStyle(document.documentElement);
    const textColor = documentStyle.getPropertyValue('--text-color');
    const textColorSecondary = documentStyle.getPropertyValue('--text-color-secondary');
    const surfaceBorder = documentStyle.getPropertyValue('--surface-border');

    let startDate = this.datePipe.transform(this.rangeDates[0], 'yyyy-MM-ddTHH:mm:ssZ');
    if (startDate != null) {
      let endDate = this.datePipe.transform(this.rangeDates[this.rangeDates.length - 1], 'yyyy-MM-ddTHH:mm:ssZ');
      let endDatedate;
      if (endDate === null) {
        endDatedate = new Date(startDate!);
        endDatedate!.setDate(endDatedate!.getDate() + 1);
        endDate = this.datePipe.transform(endDatedate, 'yyyy-MM-ddTHH:mm:ssZ');
      }
      endDate = endDate?.split('+')[0] + 'Z';
      startDate = startDate?.split('+')[0] + 'Z';

      if (this.selectedSensors != null && this.selectedSensors.length > 0 && !average) {
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
                console.log("data : ", JSON.stringify(this.data));
              }
            },
            error => {
              console.error('Error fetching measurement data:', error);
            }
          );
        });
      } else if (average && this.rangeDates != null) {
        let dates = this.getDates(this.rangeDates[0], this.rangeDates[this.rangeDates.length - 1]);
        const datasets: { label: string; data: number[]; backgroundColor: string; borderColor: string; }[] = [];
        const labels: string[] = [];

        dates.forEach(date => {
          let dateStr = this.datePipe.transform(date, 'yyyy-MM-dd');
          if (dateStr != null) {
            this.airportService.getMeasurementDataAverageByDate(this.selectedAirport!.code, dateStr).subscribe(
              measurementDataAverage => {
                measurementDataAverage.forEach(data => {
                  let color = this.getRandomColor();
                  if (datasets.find((dataset: {
                    label: string;
                  }) => dataset.label === data.measurement) === undefined) {
                    datasets.push({
                      label: data.measurement,
                      data: [data.value],
                      backgroundColor: color,
                      borderColor: color,
                    });
                  } else {
                    this.data.datasets.find((dataset: {
                      label: string;
                    }) => dataset.label === data.measurement).data.push(data.value);
                  }
                });
                labels.push(dateStr!);
                if (labels.length === dates.length) {
                  this.data = {
                    labels: labels,
                    datasets: datasets,
                  };
                }
              }
            );
          }
        });
      }
    }
    if (average){
      this.options = {
        maintainAspectRatio: false,
        aspectRatio: 0.8,
        plugins: {
          legend: {
            labels: {
              color: textColor,
            },
          },
        },
        scales: {
          x: {
            ticks: {
              color: textColorSecondary,
              font: {
                weight: 500,
              },
            },
            grid: {
              color: surfaceBorder,
              drawBorder: false,
            },
          },
          y: {
            ticks: {
              color: textColorSecondary,
            },
            grid: {
              color: surfaceBorder,
              drawBorder: false,
            },
          },
        },
      };
    }else {
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
  }

  getRandomColor() {
    const letters = '0123456789ABCDEF';
    let color = '#';
    for (let i = 0; i < 6; i++) {
      color += letters[Math.floor(Math.random() * 16)];
    }
    return color;
  }

  getDates(startDate: Date, endDate: Date): Date[] {
    const dates: Date[] = [];

    if (endDate !== null) {
      while (startDate <= endDate) {
        dates.push(new Date(startDate));
        startDate.setDate(startDate.getDate() + 1);
      }
      return dates
    }
    return [startDate];
  }
}
