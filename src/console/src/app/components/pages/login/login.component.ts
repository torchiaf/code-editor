import { Component, OnDestroy } from '@angular/core';
import { Router } from '@angular/router';
import { TranslateService } from '@ngx-translate/core';
import { Subject } from 'rxjs';
import { take } from 'rxjs/operators';
import { Route } from 'src/app/app-routing.module';
import { AuthService } from 'src/app/services/auth.service';

@Component({
  selector: 'app-login',
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.scss']
})
export class LoginComponent implements OnDestroy {

  private destroy$: Subject<void> = new Subject<void>();

  username = '';
  password = '';
  loginError?: string;
  inProgressLogin = false;

  constructor(
    private router: Router,
    private authService: AuthService,
    public translate: TranslateService
  ) {
  }

  ngOnDestroy(): void {
    this.destroy$.next();
    this.destroy$.complete();
  }

  login(): void {
    this.loginError = undefined;
    this.inProgressLogin = true;
    this.authService.login(this.username, this.password)
      .pipe(take(1))
      .subscribe({
        next: () => {
          this.loginError = undefined;
          this.router.navigateByUrl(Route.home);
        },
        error: (err) => {
          this.loginError = err.msg;
          this.inProgressLogin = false;
        },
        complete: () => this.inProgressLogin = false
      });
  }

}
