import { ComponentFixture, TestBed } from '@angular/core/testing';

import { CreateViewDialogComponent } from './create-view-dialog.component';

describe('CreateViewDialogComponent', () => {
  let component: CreateViewDialogComponent;
  let fixture: ComponentFixture<CreateViewDialogComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ CreateViewDialogComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(CreateViewDialogComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
