import { Component, EventEmitter, Input, OnDestroy, OnInit, Output } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { BehaviorSubject, Subject, catchError, combineLatest, debounceTime, from, map, switchMap, takeUntil } from 'rxjs';
import { Extension, ViewCreate } from 'src/app/models/view';
import { RestClientService } from 'src/app/services/rest-client.service';
import { ErrorDialogComponent } from '../../dialogs/error-dialog/error-dialog.component';

@Component({
  selector: 'app-view-create-form',
  templateUrl: './view-create-form.component.html',
  styleUrls: ['./view-create-form.component.scss'],
})
export class ViewCreateFormComponent implements OnInit, OnDestroy {

  @Input() data!: any;
  @Output() done = new EventEmitter<boolean | ViewCreate>();

  readonly destroyed$ = new Subject<boolean>();

  readonly accountChange$ = new BehaviorSubject<string | null>('torchiaf');

  readonly repositoryChange$ = new BehaviorSubject<string | null>('code-editor');

  readonly repositories$ = this.accountChange$.pipe(
    debounceTime(600),
    switchMap((account) => from(this.restClient.api.getRepos(account || '')).pipe(
      catchError(() => {
        this.cleanOnRepoError();
        return [];
      }),
      map((repos: any[]) => repos.map((r) => r.name))
    )),
    takeUntil(this.destroyed$));

  readonly branches$ = combineLatest([
    this.accountChange$,
    this.repositoryChange$
  ]).pipe(
    debounceTime(600),
    switchMap(([account, repo]) => from(this.restClient.api.getBranches(account || '', repo || '')).pipe(
      catchError(() => {
        this.cleanOnBranchError();
        return [];
      }),
      map((branches: any[]) => branches.map((r) => r.name)))),
      takeUntil(this.destroyed$));

  repositoryInfo = true;

  initRepo = false;
  initBranch = false;

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
        repo: 'code-editor',
        branch: 'main',
        commit: ''
      }
    }
  };

  constructor(
    private restClient: RestClientService,
    public dialog: MatDialog,
  ) {
  }

  async ngOnInit() {
    this.repositories$.subscribe((repos) => {
      this.repositories = repos;
      if (this.initRepo) {
        this.view.repo.git.repo = null;
        this.view.repo.git.branch = null;
      }
      this.initRepo = true;
    });
    this.branches$.subscribe((branches) => {
      this.branches = branches;
      if (this.initBranch) {
        this.view.repo.git.branch = null;
      }
      this.initBranch = true;
    });
  }



  ngOnDestroy(): void {
    this.accountChange$.complete();
    this.repositoryChange$.complete();
    this.destroyed$.next(true);
  }

  private cleanOnRepoError() {
    this.repositories = [];
    this.branches = [];
    this.view.repo.git.repo = null;
    this.view.repo.git.branch = null;
  }

  private cleanOnBranchError() {
    this.branches = [];
    this.view.repo.git.branch = null;
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
    try {
      // TODO fix model
      (this.view as any)['vscode-settings'] = JSON.parse(this.view.general.vscodeSettings || '{}');
      this.view.general.vscodeSettings = undefined;
    } catch (error) {
      this.dialog.open(ErrorDialogComponent, {
        width: '300px',
        height: '150px',
        data: { err: 'VSCode Settings, invalid json file' },
      });
      return;
    }

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
