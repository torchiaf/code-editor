import { Component, EventEmitter, Input, OnDestroy, OnInit, Output } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { BehaviorSubject, Subject, catchError, combineLatest, debounceTime, from, map, switchMap, takeUntil } from 'rxjs';
import { Extension, ViewCreate } from 'src/app/models/view';
import { RestClientService } from 'src/app/services/rest-client.service';
import { TranslateService } from '@ngx-translate/core';
import { FormControl, Validators } from '@angular/forms';
import { GitSelectStateMatcher, InputStateMatcher } from 'src/app/validations/state-matcher';
import { CustomValidators } from 'src/app/validations/validators';
import { ErrorDialogComponent } from '../../dialogs/error-dialog/error-dialog.component';

@Component({
  selector: 'app-view-create-form',
  templateUrl: './view-create-form.component.html',
  styleUrls: ['./view-create-form.component.scss'],
})
export class ViewCreateFormComponent implements OnInit, OnDestroy {

  nameFormControl = new FormControl('', [Validators.required]);
  emailFormControl = new FormControl('', [Validators.email]);
  gitAccountFormControl = new FormControl('', [Validators.required]);
  gitRepoFormControl = new FormControl('', [Validators.required]);
  gitBranchFormControl = new FormControl('', [Validators.required]);
  vscodeSettingsFormControl = new FormControl('', [CustomValidators.jsonValidator()]);

  inputMatcher = new InputStateMatcher();
  gitSelectMatcher = new GitSelectStateMatcher();

  @Input() data!: any;
  @Output() done = new EventEmitter<ViewCreate | null>();

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

  types = ['GitHub'];
  repositories: string[] = [];
  branches: string[] = [];

  // TODO hardcoded
  extensions: Extension[] = [{
    id: 'hoovercj.vscode-power-mode',
    settings: {
      'powermode.enabled': true
    },
    name: 'Power Mode',
  }];

  view: ViewCreate = {
    general: {
      name: '',
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
        type: 'GitHub',
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
    public translate: TranslateService,
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

  public validate(): boolean {
    return this.nameFormControl.valid &&
      this.emailFormControl.valid &&
      this.vscodeSettingsFormControl.valid &&
      (!this.repositoryInfo || (this.gitAccountFormControl.valid && this.gitRepoFormControl.valid && this.gitBranchFormControl.valid));
  }

  public save() {
    if (!this.validate()) {
      return;
    }

    try {
      this.view.general.vscodeSettings = JSON.parse(this.view.general.vscodeSettings || '{}');
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
      (this.view.repo as any) = undefined;
    }

    this.done.emit(this.view);
  }

  public cancel() {
    this.done.emit(null);
  }

}
