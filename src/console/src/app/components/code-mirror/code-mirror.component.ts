import { Component, EventEmitter, Input, OnInit, Output } from '@angular/core';

const defaultOptions: any = {
  mode: 'text/plain',
  indentWithTabs: true,
  smartIndent: true,
  lineNumbers: true,
  lineWrapping: false,
  extraKeys: { 'Ctrl-Space': 'autocomplete' },
  autoCloseBrackets: true,
  matchBrackets: true,
  lint: true,
  // readOnly: 'nocursor'
};

@Component({
  selector: 'app-code-mirror',
  templateUrl: './code-mirror.component.html',
  styleUrls: ['./code-mirror.component.scss']
})
export class CodeMirrorComponent implements OnInit {

  _options = {};

  @Input() options: any;

  @Input() data!: string;
  @Output() dataChange = new EventEmitter<string>();

  ngOnInit(): void {
    this._options = this.options ? {
      ...defaultOptions,
      ...this.options,
    } : defaultOptions;
  }
}
