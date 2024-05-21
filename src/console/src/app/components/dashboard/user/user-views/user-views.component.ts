import { ChangeDetectorRef, Component, OnDestroy, OnInit } from '@angular/core';
import { CookieService } from 'ngx-cookie';
import { FormControl } from '@angular/forms';
import { MatTableDataSource } from '@angular/material/table';
import { Subject, startWith } from 'rxjs';
import { View } from 'src/app/models/view';
import { AuthService } from 'src/app/services/auth.service';
import { RestClientService } from 'src/app/services/rest-client.service';
import { environment } from 'src/environments/environment';

type Row = View;

@Component({
  selector: 'app-user-views',
  templateUrl: './user-views.component.html',
  styleUrls: ['./user-views.component.scss']
})
export class UserViewsComponent implements OnInit, OnDestroy {

  selectedTab = new FormControl(0);

  readonly tableRefresh$: Subject<void> = new Subject<void>();

  displayedColumns = ['Id', 'Path', 'VScodeSettings', 'GoTo'];
  dataSource: MatTableDataSource<Row> = new MatTableDataSource();

  constructor(
    private cd: ChangeDetectorRef,
    private restClient: RestClientService,
    public authService: AuthService,
    private cookieService: CookieService,
  ) { }

  ngOnInit(): void {
    this.tableRefresh$.pipe(startWith(null)).subscribe(async () => {

      const views = await this.restClient.api.getViews();

      console.log(this.authService.loggedUser)

      const rows: Row[] = views.filter((v) => v.UserId === this.authService.loggedUser?.Id).sort((a, b) => a.Id.localeCompare(b.Id));

      this.dataSource = new MatTableDataSource(rows);
    });
  }

  ngOnDestroy(): void {
    this.tableRefresh$.complete();
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
