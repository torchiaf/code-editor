import { Component, EventEmitter, Input, OnInit, Output } from '@angular/core';
import { Extension, ViewCreate } from 'src/app/models/view';

@Component({
  selector: 'app-view-create-form',
  templateUrl: './view-create-form.component.html',
  styleUrls: ['./view-create-form.component.scss']
})
export class ViewCreateFormComponent {

  @Input() data!: any;
  @Output() done = new EventEmitter<boolean | ViewCreate>();

  // TODO hardcoded
  extensions: Extension[] = [{
    id: 'hoovercj.vscode-power-mode',
    settings: {
      'powermode.enabled': true
    },
    name: 'Power Mode',
  }];

  // TODO hardcoded
  types: string[] = ['gitHub'];
  accounts: string[] = ['torchiaf'];
  repositories: string[] = ['code-editor'];
  branches: string[] = ['develop', 'main'];

  view: ViewCreate = {
    general: {
      git: {
        name: 'Foo Bar',
        email: 'foo@gmail.com'
      },
      extensions: [],
      vscodeSettings: '',
    },
    repo: {
      git: {
        type: 'gitHub',
        org:'torchiaf',
        repo: 'code-editor',
        branch: 'main',
        commit: ''
      }
    }
  }

  constructor(
  ) {
  }

  save() {
    // TODO fix model
    (this.view as any)['vscode-settings'] = JSON.parse(this.view.general.vscodeSettings || '{}');
    this.view.general.vscodeSettings = undefined;

    this.view.repo.git.commit = this.view.repo.git.branch;

    this.done.emit(this.view);
  }

  cancel() {
    this.done.emit(false);
  }

}
