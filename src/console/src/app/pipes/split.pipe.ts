import { Pipe, PipeTransform } from '@angular/core';

@Pipe({
  name: 'split'
})
export class SplitPipe implements PipeTransform {

  transform(value: string, ...args: string[]): Array<string> {
    return value.split(args[0]);
  }

}
