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
});

/**
 * アイコンファイルの最大サイズ(1MB)
 */
export const MAX_ICON_FILE_SIZE = 1_000_000;
