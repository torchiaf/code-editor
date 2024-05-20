import { Component, Inject } from '@angular/core';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';
import { Extension, ViewCreate } from 'src/app/models/view';

interface ViewCreateDialogData {
  view: ViewCreate;
}

@Component({
  selector: 'app-create-view-dialog',
  templateUrl: './create-view-dialog.component.html',
  styleUrls: ['./create-view-dialog.component.scss']
})
export class CreateViewDialogComponent {

  // TODO hardcoded
  extensions: Extension[] = [{
    id: 'hoovercj.vscode-power-mode',
    settings: {
      'powermode.enabled': true
    },
    name: 'Power Mode',
  }];

  constructor(
    public dialogRef: MatDialogRef<CreateViewDialogComponent>,
    @Inject(MAT_DIALOG_DATA) public data: ViewCreateDialogData,
  ) {
  }

  save() {
    this.dialogRef.close();
  }
}
