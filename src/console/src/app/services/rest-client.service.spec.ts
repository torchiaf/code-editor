import { TestBed } from '@angular/core/testing';

import { RestClientService } from './rest-client.service';

describe('RestClientService', () => {
  let service: RestClientService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(RestClientService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
