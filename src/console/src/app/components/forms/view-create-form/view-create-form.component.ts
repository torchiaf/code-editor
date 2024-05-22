import { Component, EventEmitter, Input, Output } from '@angular/core';
import { Extension, ViewCreate } from 'src/app/models/view';

@Component({
  selector: 'app-view-create-form',
  templateUrl: './view-create-form.component.html',
  styleUrls: ['./view-create-form.component.scss'],
})
export class ViewCreateFormComponent {

  @Input() data!: any;
  @Output() done = new EventEmitter<boolean | ViewCreate>();

  repositoryInfo = true;

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

  view: any = {
    general: {
      git: {
        name: 'Foo Bar',
        email: 'foo@gmail.com'
      },
      extensions: [],
      vscodeSettings: '',
      sshKey: '',
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
  };

  public async sshFileUpload(files: File[]) {
    let sshKey = '';

    if (files?.length > 0) {
      const text = await files[0]?.text();

      sshKey = btoa(text);
    } else {
      sshKey = '';
    }

    this.view.general.sshKey = sshKey;
  }

  public save() {
    // TODO fix model
    (this.view as any)['vscode-settings'] = JSON.parse(this.view.general.vscodeSettings || '{}');
    this.view.general.vscodeSettings = undefined;

    if (this.repositoryInfo) {
      if (this.view.repo) {
        this.view.repo.git.commit = this.view.repo.git.branch;
      }
    } else {
      this.view.repo = undefined;
    }

    this.done.emit(this.view);
  }

  public cancel() {
    this.done.emit(false);
  }

}
