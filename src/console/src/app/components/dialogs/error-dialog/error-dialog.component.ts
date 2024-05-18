import { Component, Inject } from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { Error } from '../../../models/message.api';

interface ErrorComponentData {
  err: {
    error: Error
  };
}

@Component({
  selector: 'app-error-dialog',
  templateUrl: './error-dialog.component.html',
  styleUrls: ['./error-dialog.component.scss']
})
export class ErrorDialogComponent {

  constructor(
    public dialogRef: MatDialogRef<ErrorDialogComponent>,
    @Inject(MAT_DIALOG_DATA) public data: ErrorComponentData,
  ) {
    console.log(data.err);
  }

  reject(): void {
    this.dialogRef.close();
  }
}
