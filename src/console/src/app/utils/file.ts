import { saveAs as importedSaveAs } from 'file-saver';

export function downloadJson(name: string, data: any) {
  const blob = new Blob([data], { type: 'application/json' });
  importedSaveAs(blob, name);
}
