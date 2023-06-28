/**
 * 画面パス
 */
export const PAGE_PATHS = Object.freeze({
  HOME: "/",
  LOGIN: "/login",
  JOIN: "/join",
  CIRCLE_REGISTER: "/circle-register",
});

/**
 * APIパス
 */
export const API_PATHS = Object.freeze({
  USERS: "/users",
  USERS_ME: "/users/me",
  SESSIONS: "/sessions",
  CIRCLES: "/circles",
});

/**
 * メッセージ
 */
export const MESSAGES = Object.freeze({
  OVER_MAX_ICON_FILE_SIZE: "1MB以下のファイルを選択してください",
  NOT_ALLOWED_ICON_FILE_FORMAT: "jpgまたはpngファイルのみ選択できます",
  FAILED_TO_LOAD_FILE: "ファイルの読み込みに失敗しました",
});

/**
 * アイコンファイルの最大サイズ(1MB)
 */
export const MAX_ICON_FILE_SIZE = 1_000_000;

/**
 * マジックナンバー
 */
export const MAGIC_NUMBERS = Object.freeze({
  PNG: Object.freeze([0x89, 0x50, 0x4e, 0x47]),
  JPG: Object.freeze([0xff, 0xd8]),
});
