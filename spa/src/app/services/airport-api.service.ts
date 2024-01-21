import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';

interface Sensor {
  airport: string;
  measurement: string;
}

interface MeasurementData {
  airport: string;
  time: string;
  measurement: string;
  value: number;
}

@Injectable({
  providedIn: 'root',
})
export class AirportApiService {
  private apiUrl = 'http://localhost:8080';

  constructor(private http: HttpClient) {}

  getAirports(): Observable<any> {
    return this.http.get(this.apiUrl + '/airports');
  }

  getSensorsByAirportIATA(iata: string): Observable<Sensor[]> {
    return this.http.get<Sensor[]>(`${this.apiUrl}/airport/${iata}/sensors`);
  }

  getMeasurementDataByDateRangeAndType(iata: string, start: string, end: string, measurementType: string): Observable<MeasurementData[]> {
    return this.http.get<MeasurementData[]>(`${this.apiUrl}/airport/${iata}/data/range/${start}/${end}/${measurementType}`);
  }
}
