export const consts = {
  /** 有効なユーザ名 */
  ValidUsername: "testname",
  /** 有効なメールアドレス */
  ValidEmail: "foo@example.com",
  /** 有効なパスワード */
  ValidPassword: "password",
  /** 半角スペース */
  HalfWidthSpace: " ",
  /** 全角スペース */
  FullWidthSpace: "　",
  /** 半角記号 */
  HalfWidthSymbol: "`~!@#$%^&*()-_=+[]{}\\|;:'\",./<>?",
  /** 全角英字(A) */
  FullWidthA: "Ａ",
  /** 全角数字(1) */
  FullWidth1: "１",
};

export const lengths = {
  /** メールアドレスの最大桁数 */
  EmailMaxLength: 254,
  /** パスワードの最小桁数 */
  PasswordMinLength: 8,
  /** パスワードの最大桁数 */
  PasswordMaxLength: 128,
  /** ユーザ名の最大桁数 */
  UserNameMaxLength: 45,
};

export const paths = {
  /** トップページ */
  Root: "/",
  /** ログイン画面 */
  Login: "/login",
  /** アカウント作成画面 */
  Join: "/join",
};
