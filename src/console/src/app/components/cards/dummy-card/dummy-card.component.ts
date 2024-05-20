import { Component, Input } from '@angular/core';

@Component({
  selector: 'app-dummy-card',
  templateUrl: './dummy-card.component.html',
  styleUrls: ['./dummy-card.component.scss']
})
export class DummyCardComponent {
  @Input() title = 'Dummy Card';
}
