import { ComponentFixture, TestBed } from '@angular/core/testing';

import { UserViewsComponent } from './user-views.component';

describe('UserViewsComponent', () => {
  let component: UserViewsComponent;
  let fixture: ComponentFixture<UserViewsComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ UserViewsComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(UserViewsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
