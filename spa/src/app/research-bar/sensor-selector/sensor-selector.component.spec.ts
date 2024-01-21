import { ComponentFixture, TestBed } from '@angular/core/testing';

import { SensorSelectorComponent } from './sensor-selector.component';

describe('SensorSelectorComponent', () => {
  let component: SensorSelectorComponent;
  let fixture: ComponentFixture<SensorSelectorComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [SensorSelectorComponent]
    })
    .compileComponents();
    
    fixture = TestBed.createComponent(SensorSelectorComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
