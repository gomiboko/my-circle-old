import Vue from "vue";

export class Message {
  constructor(public messageType: MessageType, public message: string) {}
}

/**
 * メッセージ種別
 */
export enum MessageType {
  Success = "success",
  Info = "info",
  Warn = "warning",
  Error = "error",
}

/**
 * メッセージ表示のカスタムイベント名
 */
export const MSG_EVENT = "msg";

/**
 * エラーを表示する
 * @param vue Vueオブジェクト
 * @param error エラーオブジェクト
 */
// eslint-disable-next-line @typescript-eslint/explicit-module-boundary-types, @typescript-eslint/no-explicit-any
export function showError(vue: Vue, error: any): void {
  if (vue.$http.isAxiosError(error) && error.response && error.response.data) {
    const msg = new Message(MessageType.Error, error.response.data.message);
    vue.$emit(MSG_EVENT, msg);
  } else {
    const msg = new Message(MessageType.Error, `予期せぬエラー(${error})`);
    vue.$emit(MSG_EVENT, msg);
  }
}
