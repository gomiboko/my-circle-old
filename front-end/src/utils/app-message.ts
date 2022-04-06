/**
 * アプリケーションメッセージクラス。
 * 画面上部に表示するメッセージ。
 */
export class AppMessage {
  private _md!: number;
  private _lg!: number;
  private _xl!: number;

  /**
   * 指定したメッセージ種別、メッセージ、表示サイズでAppMessageを生成する
   * @param type メッセージ種別
   * @param message メッセージ
   * @param size 表示サイズ
   */
  constructor(public type: AppMessageType, public message: string, size: AppMessageSize = AppMessageSize.Medium) {
    this.setSize(size);
  }

  /**
   * 画面サイズが md のときの表示サイズ
   */
  get md(): number {
    return this._md;
  }

  /**
   * 画面サイズが lg のときの表示サイズ
   */
  get lg(): number {
    return this._lg;
  }

  /**
   * 画面サイズが xl のときの表示サイズ
   */
  get xl(): number {
    return this._xl;
  }

  /**
   * メッセージの表示サイズを設定する
   * @param size 表示サイズ
   */
  public setSize(size: AppMessageSize): void {
    switch (size) {
      case AppMessageSize.Small:
        this._md = 4;
        this._lg = 3;
        this._xl = 2;
        break;
      case AppMessageSize.Medium:
        this._md = 8;
        this._lg = 6;
        this._xl = 4;
        break;
      case AppMessageSize.Large:
        this._md = 12;
        this._lg = 9;
        this._xl = 6;
        break;
    }
  }
}

/**
 * メッセージ種別
 */
export enum AppMessageType {
  Success = "success",
  Info = "info",
  Warn = "warning",
  Error = "error",
}

/**
 * メッセージ表示サイズ
 */
export enum AppMessageSize {
  Small,
  Medium,
  Large,
}
