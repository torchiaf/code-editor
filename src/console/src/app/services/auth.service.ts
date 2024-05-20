import { HttpResponse } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Router } from '@angular/router';
import { Observable, throwError, from } from 'rxjs';
import { catchError, map } from 'rxjs/operators';
import { environment } from 'src/environments/environment';
import jwt_decode from 'jwt-decode';
import * as _ from 'lodash';
import { RestClientService } from './rest-client.service';
import { Route } from '../app-routing.module';
import { UserDetails } from '../models/user';

interface LoggedUser extends Omit<UserDetails, 'roleId'> {
  expirationDate: Date;
}

@Injectable({
  providedIn: 'root'
})
export class AuthService {

  public loggedUser?: LoggedUser;

  private toLoggedUser(jwtToken: string): LoggedUser {
    const user: UserDetails & { exp: number } = jwt_decode(jwtToken);
    const expirationDate = new Date(Math.floor(new Date().getTime() + user.exp));
    return {
      Id: user.Id,
      Name: user.username,
      role: user.role,
      expirationDate
    };
  }

  constructor(
    private router: Router,
    private restClient: RestClientService,
  ) {
    this.authenticate();
  }

  public get isAuthenticated() {
    return !_.isEmpty(this.loggedUser);
  }

  private authenticate(): void {
    const jwtToken = this.getJwtToken();
    this.saveUser(jwtToken);
  }

  public getJwtToken(): string | null {
    return localStorage.getItem(environment.jwtToken);
  }

  private storeJwtToken(jwtToken: string): void {
    localStorage.setItem(environment.jwtToken, jwtToken);
  }

  private discardToken(): void {
    localStorage.removeItem(environment.jwtToken);
    this.loggedUser = undefined;
  }

  private saveUser(jwtToken: string | null): void {
    if (jwtToken) {
      this.loggedUser = this.toLoggedUser(jwtToken);
      // console.info('Logged user:', this.loggedUser);
    }
  }

  public isLoginExpired(): boolean {
    if (!this.loggedUser) {
      return true;
    }
    return Date.now() > this.loggedUser.expirationDate.getTime();
  }

  public login(username: string, password: string): Observable<void> {
    return from(this.restClient.api.login(username, password)).pipe(
      catchError((err: HttpResponse<Error>) => {
        console.error(err);
        this.discardToken();
        return throwError(() => ({
          msg: err.status === 0 ? 'ERROR_LOGIN_CONN' : 'ERROR_LOGIN_INVALID'
        }));
      }),
      map((res: any) => {
        const token = res[environment.jwtToken];
        this.storeJwtToken(token);
        this.saveUser(token);
      }),
    );
  }

  public logout() {
    this.discardToken();
    this.router.navigateByUrl(Route.login);
  }

}
