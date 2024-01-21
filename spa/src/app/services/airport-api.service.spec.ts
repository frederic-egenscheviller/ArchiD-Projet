import { TestBed } from '@angular/core/testing';

import { AirportApiService } from './airport-api.service';

describe('AirportApiService', () => {
  let service: AirportApiService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(AirportApiService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
