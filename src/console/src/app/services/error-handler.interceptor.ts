import { Injectable } from '@angular/core';
import {
  HttpRequest,
  HttpHandler,
  HttpEvent,
  HttpInterceptor,
} from '@angular/common/http';
import { Observable } from 'rxjs';
import { tap } from 'rxjs/operators';
import { ErrorHandlerService } from './error-handler.service';

@Injectable({
  providedIn: 'root'
})
export class ErrorHandlerInterceptor implements HttpInterceptor {

  constructor(
    private errorHandlerService: ErrorHandlerService
  ) {
  }

  intercept(request: HttpRequest<unknown>, next: HttpHandler): Observable<HttpEvent<unknown>> {
    if (request.url.includes('//api.github.com/')) {
      return next.handle(request).pipe(
        tap({
          error: (err) => {
            if (err.status === 403) {
              return this.errorHandlerService.httpResponseError$.next('GitHub API rate limit exceeded' as any);
            }
            return null;
          },
        })
      );
    }

    return next.handle(request).pipe(
      tap({
        error: (err) => this.errorHandlerService.httpResponseError$.next(err),
      }),
    );
  }
}
