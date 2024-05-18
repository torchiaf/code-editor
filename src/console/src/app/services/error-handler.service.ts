import { HttpResponse } from '@angular/common/http';
import { Injectable, OnDestroy } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { BehaviorSubject } from 'rxjs';
import { debounceTime, filter } from 'rxjs/operators';
import { Route } from '../app-routing.module';
import { ErrorDialogComponent } from '../components/dialogs/error-dialog/error-dialog.component';

@Injectable({
  providedIn: 'root'
})
export class ErrorHandlerService implements OnDestroy {

  public httpResponseError$ = new BehaviorSubject<HttpResponse<any> | undefined>(undefined);

  constructor(
    public dialog: MatDialog,
  ) {
    this.httpResponseError$.asObservable().pipe(
      filter((f) => f !== undefined),
      filter((f) => !f?.url?.includes(Route.login)),
      debounceTime(200),
    ).subscribe((err) => {
      this.dialog.open(ErrorDialogComponent, {
        width: '300px',
        height: '150px',
        data: { err },
      });
    });
  }

  ngOnDestroy(): void {
    this.httpResponseError$.complete();
  }

}
