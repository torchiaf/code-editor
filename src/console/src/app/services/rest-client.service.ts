import { HttpClient, HttpContext, HttpHeaders, HttpParams } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { lastValueFrom } from 'rxjs';
import { map } from 'rxjs/operators';
import { environment } from 'src/environments/environment';
import { Login } from '../models/login';
import { UserDetails } from '../models/user';
import { View, ViewCreateGeneral, ViewCreateRepo } from '../models/view';

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

  private http = (url: string, callback = <Z>({ data }: { data: Z }) => data) => ({
    get: <T>(
      path: string,
      options?: HttpGetOptionParams,
    ) => lastValueFrom(this.httpClient.get<ResponseApi<T>>(`${url}/${path}`, options).pipe(map(callback))),
    post: <T>(
      path: string,
      body: any | null,
      options?: HttpGetOptionParams
    ) => lastValueFrom(this.httpClient.post<ResponseApi<T>>(`${url}/${path}`, body, options).pipe(map(callback))),
    put: <T>(
      path: string,
      body: any | null,
      options?: HttpGetOptionParams
    ) => lastValueFrom(this.httpClient.put<ResponseApi<T>>(`${url}/${path}`, body, options).pipe(map(callback))),
    delete: <T>(
      path: string,
      options?: HttpGetOptionParams
    ) => lastValueFrom(this.httpClient.delete<ResponseApi<T>>(`${url}/${path}`, options).pipe(map(callback))),
  });

  private codeEditorApi = this.http(`${environment.protocol}://${window.location.hostname}${environment.restPath}/${environment.apiVersion}`);
  private gitHubApi = this.http('https://api.github.com', (x: any) => x);

  constructor(
    private httpClient: HttpClient,
  ) {
  }

  public api = {

    ping: () => this.codeEditorApi.get<UserDetails>('ping'),

    /** Login */

    login: (username: string, password: string) => this.codeEditorApi.post<Login>('login', { username, password }),

    /** User */

    getUsers: () => this.codeEditorApi.get<Array<UserDetails>>('users'),

    getUser: (name: string) => this.codeEditorApi.get<UserDetails>(`user/${name}`),

    /** Views */

    getViews: () => this.codeEditorApi.get<Array<View>>('views'),

    getView: (viewId: string) => this.codeEditorApi.get<View>(`views/${viewId}`),

    createView: (username: string, viewGeneral: ViewCreateGeneral) => this.codeEditorApi.post<{ viewId: string }>(`views?username=${username}`, viewGeneral),

    updateView: (viewId: string, viewRepo: ViewCreateRepo) => this.codeEditorApi.put<void>(`views/${viewId}`, viewRepo),

    deleteView: (viewId: string) => this.codeEditorApi.delete<void>(`views/${viewId}`),

    /** GitHub */

    getRepos: (org: string) => this.gitHubApi.get<any[]>(`users/${org}/repos`),

    getBranches: (org: string, repo: string) => this.gitHubApi.get<any[]>(`repos/${org}/${repo}/branches`),

  } as const;


}
