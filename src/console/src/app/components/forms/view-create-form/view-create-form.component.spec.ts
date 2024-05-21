import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ViewCreateFormComponent } from './view-create-form.component';

describe('ViewCreateFormComponent', () => {
  let component: ViewCreateFormComponent;
  let fixture: ComponentFixture<ViewCreateFormComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ ViewCreateFormComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(ViewCreateFormComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
