import { Component, OnDestroy, OnInit } from '@angular/core';
import { NavigationEnd, Router } from '@angular/router';
import { TranslateService } from '@ngx-translate/core';
import { filter, lastValueFrom, Subject, takeUntil } from 'rxjs';
import { environment } from 'src/environments/environment';
import { AuthService } from './services/auth.service';
import { MatDialog } from '@angular/material/dialog';
import { Route } from './app-routing.module';
import { ProfileDialogComponent } from './components/dialogs/profile-dialog/profile-dialog.component';
import { MatIconRegistry } from '@angular/material/icon';
import { DomSanitizer } from '@angular/platform-browser';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent implements OnInit, OnDestroy {

  private destroy$: Subject<void> = new Subject<void>();

  public currentUrl = Route.home;
  readonly Route = Route;
  readonly languages: Array<string> = environment.languages;

  constructor(
    public dialog: MatDialog,
    public authService: AuthService,
    public translate: TranslateService,
    private router: Router,
    private matIconRegistry: MatIconRegistry,
    private domSanitizer: DomSanitizer,
  ) {
    translate.addLangs(this.languages);
    translate.use(environment.defaultLanguage);

    this.matIconRegistry.addSvgIcon(
      'code-editor',
      this.domSanitizer.bypassSecurityTrustResourceUrl('../assets/logo-black.svg'),
    );
  }

  ngOnInit(): void {
    this.router.events.pipe(
      filter((event) => event instanceof NavigationEnd),
      takeUntil(this.destroy$)
    ).subscribe((r: any) => {
      this.currentUrl = r.url;
    });
  }

  ngOnDestroy(): void {
    this.destroy$.next();
    this.destroy$.complete();
  }

  switchTranslations(lang: string): void {
    this.translate.use(lang);
  }

  gitHubHome(): void {
    window.open(environment.gitHubHome, '_blank');
  }

  goHome(): void {
    this.router.navigateByUrl('/');
  }

  async showProfile() {
    const dialogRef = this.dialog.open(ProfileDialogComponent, {
      width: '200px',
      height: '250px',
      data: { user: this.authService.loggedUser },
    });
    const res = await lastValueFrom(dialogRef.afterClosed());
    if (res) {
      console.log('Profile Card');
    }
  }

  logout(): void {
    this.authService.logout();
  }

}
