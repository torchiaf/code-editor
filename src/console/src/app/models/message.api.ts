interface Msg {
  message: string;
}

export interface Success extends Msg {
  data: any;
}

export interface Error extends Msg {
  error: string;
}
