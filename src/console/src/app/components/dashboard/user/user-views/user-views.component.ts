import { Component, OnDestroy, OnInit } from '@angular/core';
import { CookieService } from 'ngx-cookie';
import { FormControl } from '@angular/forms';
import { Subject, lastValueFrom, startWith } from 'rxjs';
import { View, ViewCreate } from 'src/app/models/view';
import { AuthService } from 'src/app/services/auth.service';
import { RestClientService } from 'src/app/services/rest-client.service';
import { environment } from 'src/environments/environment';
import { MatDialog } from '@angular/material/dialog';
import { ConfirmDialogComponent } from 'src/app/components/dialogs/confirm-dialog/confirm-dialog.component';

@Component({
  selector: 'app-user-views',
  templateUrl: './user-views.component.html',
  styleUrls: ['./user-views.component.scss']
})
export class UserViewsComponent implements OnInit, OnDestroy {

  selectedTab = new FormControl(0);
  createView = false;

  readonly cardRefresh$: Subject<void> = new Subject<void>();
  data: View | null = null;

  creating = false;
  deleting = false;

  constructor(
    public dialog: MatDialog,
    private restClient: RestClientService,
    public authService: AuthService,
    private cookieService: CookieService,
  ) { }

  ngOnInit(): void {
    this.cardRefresh$.pipe(startWith(null)).subscribe(async () => {
      const views = await this.restClient.api.getViews();

      const rows: View[] = views.filter((v) => v.UserId === this.authService.loggedUser?.Id);

      this.data = rows[0];
    });
  }

  ngOnDestroy(): void {
    this.cardRefresh$.complete();
  }

  public goToView(element: View) {
    this.cookieService.get('code-server-session');
    this.cookieService.put('code-server-session', element.Session, {
      path: '/',
      secure: false,
      storeUnencoded: true
    });


    const url = `${environment.protocol}://${window.location.hostname}${element.Path}?${element.Query}`;

    window.open(url, '_blank');
  }

  public goToCreateView() {
    this.createView = true;

    this.selectedTab.setValue(1);
  }

  public goToViews() {
    this.createView = false;

    this.selectedTab.setValue(0);
  }

  public async deleteView() {
    const dialogRef = this.dialog.open(ConfirmDialogComponent, {
      width: '300px',
      height: '150px',
      data: { message: 'DELETE_VIEW_CONFIRM_MSG', type: 'delete' },
    });
    const res = await lastValueFrom(dialogRef.afterClosed());
    if (res && this.data) {
      this.deleting = true;

      await this.restClient.api.deleteView(this.data.Id);

      this.deleting = false;

      this.cardRefresh$.next();
    }
  }

  public async createViewDone(res: boolean | ViewCreate) {
    if (res) {
      this.creating = true;
      this.goToViews();

      try {
        const created = await this.restClient.api.userCreateView((res as ViewCreate).general);

        const repoInfo = (res as ViewCreate).repo;
        if(repoInfo) {
          await this.restClient.api.updateView(created.viewId || '', repoInfo);
        }
      } catch (error) {
        this.creating = false;
      }

      this.creating = false;
      this.cardRefresh$.next();
    }

    this.goToViews();
  }

}
