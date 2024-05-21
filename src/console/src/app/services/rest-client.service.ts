import { HttpClient, HttpContext, HttpHeaders, HttpParams } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { lastValueFrom } from 'rxjs';
import { map } from 'rxjs/operators';
import { environment } from 'src/environments/environment';
import { Login } from '../models/login';
import { UserDetails } from '../models/user';
import { View, ViewCreate, ViewCreateGeneral, ViewCreateRepo } from '../models/view';

interface HttpGetOptionParams {
  headers?: HttpHeaders | {
    [header: string]: string | string[];
  };
  context?: HttpContext;
  observe?: 'body';
  params?: HttpParams | {
    [param: string]: string | boolean | ReadonlyArray<string | boolean>;
  };
  reportProgress?: boolean;
  responseType?: 'json';
  withCredentials?: boolean;
}

interface ResponseApi<T> {
  message: string;
  data: T;
}

@Injectable({
  providedIn: 'root'
})
export class RestClientService {

  private readonly url = `${environment.restURL}/${environment.apiVersion}`;

  // private headers: Headers;

  private http = {
    get: <T>(
      path: string,
      options?: HttpGetOptionParams
    ) => lastValueFrom(this.httpClient.get<ResponseApi<T>>(`${this.url}/${path}`, options).pipe(map(({ data }) => data))),
    post: <T>(
      path: string,
      body: any | null,
      options?: HttpGetOptionParams
    ) => lastValueFrom(this.httpClient.post<ResponseApi<T>>(`${this.url}/${path}`, body, options).pipe(map(({ data }) => data))),
    put: <T>(
      path: string,
      body: any | null,
      options?: HttpGetOptionParams
    ) => lastValueFrom(this.httpClient.put<ResponseApi<T>>(`${this.url}/${path}`, body, options).pipe(map(({ data }) => data))),
    delete: <T>(
      path: string,
      options?: HttpGetOptionParams
    ) => lastValueFrom(this.httpClient.delete<ResponseApi<T>>(`${this.url}/${path}`, options).pipe(map(({ data }) => data))),
  } as const;

  constructor(
    private httpClient: HttpClient,
  ) {
  }

  public api = {

    ping: () => this.http.get<UserDetails>('ping'),

    /** Login */

    login: (username: string, password: string) => this.http.post<Login>('login', { username, password }),

    /** User */

    getUsers: () => this.http.get<Array<UserDetails>>('users'),

    getUser: (name: string) => this.http.get<UserDetails>(`user/${name}`),

    /** Views */

    getViews: () => this.http.get<Array<View>>('views'),

    getView: (viewId: string) => this.http.get<View>(`views/${viewId}`),

    createView: (username: string, viewGeneral: ViewCreateGeneral) => this.http.post<{ viewId: string }>(`views?username=${username}`, viewGeneral),

    updateView: (viewId: string, viewRepo: ViewCreateRepo) => this.http.put<void>(`views/${viewId}`, viewRepo),

    deleteView: (viewId: string) => this.http.delete<void>(`views/${viewId}`),

  } as const;


}
