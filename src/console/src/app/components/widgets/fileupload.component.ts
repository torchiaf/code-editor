import { Component, Input, ViewChild, ElementRef, Output, EventEmitter } from "@angular/core";
import { DomSanitizer, SafeUrl } from '@angular/platform-browser';

@Component({
  selector: 'vo-mat-fileUpload',
  templateUrl: './fileupload.component.html',
  styleUrls: ['./fileupload.component.scss']
})
export class FileuploadComponent {
  @Input() mode: any
  @Input() names: any
  @Input() url: any
  @Input() method: any
  @Input() multiple: any
  @Input() disabled: any
  @Input() accept: any
  @Input() maxFileSize: any
  @Input() auto = true
  @Input() withCredentials: any
  @Input() invalidFileSizeMessageSummary: any
  @Input() invalidFileSizeMessageDetail: any
  @Input() invalidFileTypeMessageSummary: any
  @Input() invalidFileTypeMessageDetail: any
  @Input() previewWidth: any
  @Input() chooseLabel = 'Choose'
  @Input() uploadLabel = 'Upload'
  @Input() cancelLabel = 'Cance'
  @Input() customUpload: any
  @Input() showUploadButton: any
  @Input() showCancelButton: any

  @Input() dataUriPrefix: any
  @Input() deleteButtonLabel: any
  @Input() deleteButtonIcon = 'close'
  @Input() showUploadInfo: any

  @ViewChild('fileUpload') fileUpload: ElementRef = new ElementRef('')
  inputFileName: string = ''

  @Input() files: File[] = []

  @Output() done = new EventEmitter<File[]>();

  constructor(private sanitizer: DomSanitizer) {
  }

  onClick(event: any) {
    if (this.fileUpload)
      this.fileUpload.nativeElement.click()
  }

  onInput(event: any) {
  }

  onFileSelected(event: any) {
    let files = event.dataTransfer ? event.dataTransfer.files : event.target.files;
    // console.log('event::::::', event)
    for (let i = 0; i < files.length; i++) {
      let file = files[i];

      if (this.validate(file)) {
        file.objectURL = this.sanitizer.bypassSecurityTrustUrl((window.URL.createObjectURL(files[i])));

        if (!this.isMultiple()) {
          this.files = []
        }
        this.files.push(files[i]);
      }
    }

    this.done.emit(this.files);
  }

  removeFile(event: any, file: any) {
    let ix
    if (this.files && -1 !== (ix = this.files.indexOf(file))) {
      this.files.splice(ix, 1)
      this.clearInputElement()
    }

    this.done.emit(this.files);
  }

  validate(file: File) {
    for (const f of this.files) {
      if (f.name === file.name
        && f.lastModified === file.lastModified
        && f.size === f.size
        && f.type === f.type
      ) {
        return false
      }
    }
    return true
  }

  clearInputElement() {
    this.fileUpload.nativeElement.value = ''
  }


  isMultiple(): boolean {
    return this.multiple
  }
}
