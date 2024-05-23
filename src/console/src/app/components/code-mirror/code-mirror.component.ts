import { Component, Input } from '@angular/core';

@Component({
  selector: 'app-code-mirror',
  templateUrl: './code-mirror.component.html',
  styleUrls: ['./code-mirror.component.scss']
})
export class CodeMirrorComponent {

  @Input() data!: string;

  codeMirrorOptions: any = {
    mode: 'text/plain',
    indentWithTabs: true,
    smartIndent: true,
    lineNumbers: true,
    lineWrapping: false,
    extraKeys: { 'Ctrl-Space': 'autocomplete' },
    autoCloseBrackets: true,
    matchBrackets: true,
    lint: true,
    readOnly: 'nocursor'
  };

}
