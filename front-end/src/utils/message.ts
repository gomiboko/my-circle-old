export class Message {
  constructor(public messageType: MessageType, public message: string) {}
}

export const enum MessageType {
  Success = "success",
  Info = "info",
  Warn = "warning",
  Error = "error",
}

/**
 * メッセージ表示のカスタムイベント名
 */
export const MSG_EVENT = "msg";
