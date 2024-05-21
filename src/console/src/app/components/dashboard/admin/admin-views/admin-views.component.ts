import { animate, state, style, transition, trigger } from '@angular/animations';
import { ChangeDetectorRef, Component, OnDestroy, OnInit, QueryList, ViewChild, ViewChildren } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { MatSort } from '@angular/material/sort';
import { MatTable, MatTableDataSource } from '@angular/material/table';
import { lastValueFrom, Subject } from 'rxjs';
import { startWith } from 'rxjs/operators';
import { ConfirmDialogComponent } from 'src/app/components/dialogs/confirm-dialog/confirm-dialog.component';
import { UserDetails } from 'src/app/models/user';
import { View, ViewCreate } from 'src/app/models/view';
import { AuthService } from 'src/app/services/auth.service';
import { RestClientService } from 'src/app/services/rest-client.service';
import { FormControl } from '@angular/forms';
import { environment } from 'src/environments/environment';
import { CookieService } from 'ngx-cookie';

type Row = UserDetails & { Enabled: boolean, Views: View[] | MatTableDataSource<View> };

@Component({
  selector: 'app-admin-views',
  templateUrl: './admin-views.component.html',
  styleUrls: ['./admin-views.component.scss'],
  animations: [
    trigger('detailExpand', [
      state('collapsed', style({ height: '0px', minHeight: '0' })),
      state('expanded', style({ height: '*' })),
      transition('expanded <=> collapsed', animate('225ms cubic-bezier(0.4, 0.0, 0.2, 1)')),
    ]),
  ],
})
export class AdminViewsComponent implements OnInit, OnDestroy {

  @ViewChild('outerSort', { static: true }) sort: MatSort | undefined;
  @ViewChildren('innerSort') innerSort: QueryList<MatSort> | undefined;
  @ViewChildren('innerTables') innerTables: QueryList<MatTable<View>> | undefined;

  selectedTab = new FormControl(0);
  createView = false;

  readonly tableRefresh$: Subject<void> = new Subject<void>();

  dataSource: MatTableDataSource<Row> = new MatTableDataSource();

  displayedColumns = ['Id', 'Name', 'Email', 'Phone', 'Status'];
  innerDisplayedColumns = ['Id', 'Path', 'VScodeSettings', 'Delete', 'GoTo'];

  expandedElements: Record<string, boolean> = {};
  collapseDisabled = false;

  creating: string | null = null;
  deleting: string | null = null;

  createViewData: any;

  constructor(
    private cd: ChangeDetectorRef,
    public dialog: MatDialog,
    private restClient: RestClientService,
    public authService: AuthService,
    private cookieService: CookieService,
  ) {
  }

  ngOnInit(): void {
    this.tableRefresh$.pipe(startWith(null)).subscribe(async () => {

      const [views, users] = await Promise.all([
        this.restClient.api.getViews(),
        this.restClient.api.getUsers()
      ]);

      const rows: Row[] = users.map((user) => {
        const userViews = views.filter((v) => v.UserId === user.Id);

        return {
          ...user,

          // TODO hardcoded
          Email: 'foo@gmail.com',
          Phone: '123456789',

          Enabled: userViews.length > 0,
          Views: new MatTableDataSource(userViews)
        };
      }).sort((a, b) => a.Id.localeCompare(b.Id));

      this.dataSource = new MatTableDataSource(rows);

      this.expandedElements = rows.reduce((acc, r) => ({
        ...acc,
        [r.Id]: this.expandedElements[r.Id] || false,
      }), {});
      this.collapseDisabled = !rows.find((r) => (r.Views as MatTableDataSource<View>).data.length > 0);
    });
  }

  ngOnDestroy(): void {
    this.tableRefresh$.complete();
  }

  toggleRow(element: Row) {
    this.expandedElements[element.Id] = !this.expandedElements[element.Id];

    this.cd.detectChanges();
  }

  public toggleCollapse() {
    const toggle = Object.keys(this.expandedElements).some((id) => !this.expandedElements[id] && this.dataSource.data.find((f) => f.Id === id && (f.Views as MatTableDataSource<View>).data.length > 0));
    Object.keys(this.expandedElements).forEach((id) => this.expandedElements[id] = toggle);

    this.cd.detectChanges();
  }

  public goToCreateView(element: Row) {
    this.createView = true;
    this.createViewData = element;

    this.selectedTab.setValue(1);
  }

  public goToViews() {
    this.createView = false;
    this.createViewData = null;

    this.selectedTab.setValue(0);
  }

  public async createViewDone(res: boolean | ViewCreate) {
    if (res) {
      this.creating = this.createViewData.Id;

      const row = this.dataSource.data.find((row) => row.Id === this.creating);

      if (row) {
        this.goToViews();

        try {
          const created = await this.restClient.api.createView(row.Name || '', (res as ViewCreate).general);

          const repoInfo = (res as ViewCreate).repo;
          if(repoInfo) {
            await this.restClient.api.updateView(created.viewId || '', repoInfo);
          }    
        } catch (error) {
          this.creating = null;
        }

        this.creating = null;
        this.tableRefresh$.next();
      }
    }

    this.goToViews();
  }

  public async deleteView(view: Row) {
    const dialogRef = this.dialog.open(ConfirmDialogComponent, {
      width: '300px',
      height: '150px',
      data: { message: 'DELETE_VIEW_CONFIRM_MSG', type: 'delete' },
    });
    const res = await lastValueFrom(dialogRef.afterClosed());
    if (res) {
      this.deleting = view.Id;

      await this.restClient.api.deleteView(view.Id);

      this.deleting = null;

      this.tableRefresh$.next();
    }
  }

  public goToView(element: View) {
    this.cookieService.get('code-server-session');
    this.cookieService.put('code-server-session', element.Session, {
      path: '/',
      secure: false,
      storeUnencoded: true
    });

    const url = `${environment.baseUrl}${element.Path}?${element.Repo}`;

    window.open(url, '_blank');
  }
}
