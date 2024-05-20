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
import { CreateViewDialogComponent } from '../../dialogs/create-view-dialog/create-view-dialog.component';

type Row = UserDetails & { Enabled: boolean, Views: View[] | MatTableDataSource<View> };

@Component({
  selector: 'app-views',
  templateUrl: './views.component.html',
  styleUrls: ['./views.component.scss'],
  animations: [
    trigger('detailExpand', [
      state('collapsed', style({ height: '0px', minHeight: '0' })),
      state('expanded', style({ height: '*' })),
      transition('expanded <=> collapsed', animate('225ms cubic-bezier(0.4, 0.0, 0.2, 1)')),
    ]),
  ],
})
export class ViewsComponent implements OnInit, OnDestroy {

  @ViewChild('outerSort', { static: true }) sort: MatSort | undefined;
  @ViewChildren('innerSort') innerSort: QueryList<MatSort> | undefined;
  @ViewChildren('innerTables') innerTables: QueryList<MatTable<View>> | undefined;

  readonly tableRefresh$: Subject<void> = new Subject<void>();

  dataSource: MatTableDataSource<Row> = new MatTableDataSource();

  displayedColumns = ['Id', 'Name', 'Email', 'Phone', 'Status'];
  innerDisplayedColumns = ['Id', 'Path', 'VScodeSettings', 'Delete'];

  expandedElements: Record<string, boolean> = {};
  collapseDisabled = false;

  creating: string | null = null;
  deleting: string | null = null;

  constructor(
    private cd: ChangeDetectorRef,
    public dialog: MatDialog,
    private restClient: RestClientService,
    public authService: AuthService,
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

  public async addView(row: Row) {
    const dialogRef = this.dialog.open(CreateViewDialogComponent, {
      width: '700px',
      height: '600px',
      data: {view: { git: {}}},
      disableClose: true
    });
    const res: ViewCreate = await lastValueFrom(dialogRef.afterClosed());
    if (res) {
      // TODO fix model
      (res as any)['vscode-settings'] = JSON.parse(res.vscodeSettings || '{}');
      res.vscodeSettings = undefined;

      this.creating = row.Id;

      await this.restClient.api.createView(row.Name || '', res);

      this.creating = null;
      this.tableRefresh$.next();
    }
  }

  public async deleteView(view: Row) {
    const dialogRef = this.dialog.open(ConfirmDialogComponent, {
      width: '300px',
      height: '150px',
      data: { message: 'DELETE_VIEW_CONFIRM_MSG' },
    });
    const res = await lastValueFrom(dialogRef.afterClosed());
    if (res) {
      this.deleting = view.Id;

      await this.restClient.api.deleteView(view.Id);

      this.deleting = null;

      this.tableRefresh$.next();
    }
  }
}
