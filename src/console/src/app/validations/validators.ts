import { AbstractControl, ValidationErrors, ValidatorFn } from '@angular/forms';

function jsonValidator(): ValidatorFn {
  return (control: AbstractControl): ValidationErrors | null => {
    try {
      JSON.parse(control.value || '{}');
      return null;
    } catch (error) {
      return { invalidJson: { value: control.value } };
    }
  };
}

export const CustomValidators = {
  jsonValidator,
};
