import {Component, Output} from '@angular/core';

@Component({
  selector: 'app-research-bar',
  templateUrl: './research-bar.component.html',
  styleUrls: ['./research-bar.component.scss']
})
export class ResearchBarComponent {
  @Output()
  airportIsSelected = false;
}
