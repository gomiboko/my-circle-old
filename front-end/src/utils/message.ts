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
 * @param colSize エラーメッセージの表示サイズ
 */
// eslint-disable-next-line @typescript-eslint/explicit-module-boundary-types, @typescript-eslint/no-explicit-any
export function showError(vue: Vue, error: any, colSize?: number): void {
  let msg: Message;
  if (vue.$http.isAxiosError(error) && error.response && error.response.data) {
    msg = new Message(MessageType.Error, error.response.data.message);
  } else {
    msg = new Message(MessageType.Error, `予期せぬエラー(${error})`);
  }

  vue.$emit(MSG_EVENT, msg, colSize);
}
