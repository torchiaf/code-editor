import { Component, EventEmitter, Input, OnDestroy, OnInit, Output } from '@angular/core';
import { BehaviorSubject, catchError, combineLatest, debounceTime, from, map, switchMap } from 'rxjs';
import { Extension, ViewCreate } from 'src/app/models/view';
import { RestClientService } from 'src/app/services/rest-client.service';

@Component({
  selector: 'app-view-create-form',
  templateUrl: './view-create-form.component.html',
  styleUrls: ['./view-create-form.component.scss'],
})
export class ViewCreateFormComponent implements OnInit, OnDestroy {

  @Input() data!: any;
  @Output() done = new EventEmitter<boolean | ViewCreate>();

  readonly accountChange$ = new BehaviorSubject<string | null>('torchiaf');

  readonly repositoryChange$ = new BehaviorSubject<string | null>(null);

  readonly repositories$ = this.accountChange$.pipe(
    debounceTime(600),
    switchMap((account) => from(this.restClient.api.getRepos(account || '')).pipe(
      catchError(() => []),
      map((repos: any[]) => repos.map((r) => r.name))
    )));

  readonly branches$ = combineLatest([
    this.accountChange$,
    this.repositoryChange$
  ]).pipe(
    debounceTime(600),
    switchMap(([account, repo]) => from(this.restClient.api.getBranches(account || '', repo || 'code-editor')).pipe(
      catchError(() => []),
      map((branches: any[]) => branches.map((r) => r.name)))));

  repositoryInfo = true;

  types: string[] = ['gitHub'];
  repositories: any[] = [];
  branches: string[] = [];

  // TODO hardcoded
  extensions: Extension[] = [{
    id: 'hoovercj.vscode-power-mode',
    settings: {
      'powermode.enabled': true
    },
    name: 'Power Mode',
  }];

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
        repo: '',
        branch: '',
        commit: ''
      }
    }
  };

  constructor(
    private restClient: RestClientService,
  ) {
  }

  async ngOnInit() {
    this.repositories$.subscribe((repos) => {
      this.repositories = repos;
      this.view.repo.git.repo = null;
      this.view.repo.git.branch = null;
    });
    this.branches$.subscribe((branches) => {
      this.branches = branches;
      this.view.repo.git.branch = null;
    });
  }

  ngOnDestroy(): void {
    this.accountChange$.complete();
    this.repositoryChange$.complete();
  }

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

  public validateRepositoryInfo(): boolean {
    return !this.repositoryInfo || !!(
      this.view?.repo?.git?.type &&
      this.view?.repo?.git?.org &&
      this.view?.repo?.git?.repo &&
      this.view?.repo?.git?.branch
    );
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
