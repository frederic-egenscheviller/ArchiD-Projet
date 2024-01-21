import {Component, EventEmitter, OnInit, Output} from '@angular/core';
import {DropdownModule} from "primeng/dropdown";
import {FormsModule} from "@angular/forms";

interface Airport {
  code: string;
}
@Component({
  selector: 'app-airport-selector',
  standalone: true,
  imports: [
    DropdownModule,
    FormsModule
  ],
  templateUrl: './airport-selector.component.html',
  styleUrl: './airport-selector.component.scss'
})
export class AirportSelectorComponent implements OnInit{
  @Output() airportIsSelected: EventEmitter<boolean> = new EventEmitter<boolean>();

  airports: Airport[] | undefined;

  selectedAirport: Airport | undefined;

  ngOnInit() {
    this.airports = [
      { code: 'MRS' },
      { code: 'LYS' }
    ];
  }
}
