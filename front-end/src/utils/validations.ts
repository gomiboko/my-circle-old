import { ValidationRule } from "vee-validate/dist/types/types";
import { email } from "vee-validate/dist/rules";

/**
 * VeeValidate標準の検証ルールemailをカスタムした検証ルール。
 * デフォルトの日本語メッセージが不自然なため、メッセージのみ上書き。
 */
export const customEmail: ValidationRule = {
  ...email,
  message: "有効なメールアドレスではありません",
};

/**
 * パスワードの検証ルール
 */
export const password: ValidationRule = {
  validate(value) {
    return /^[a-zA-Z0-9!@#$%^&*()-_=+[\]{}\\|~;:'",.<>/?`]*$/.test(value);
  },
  message: "{_field_}に使用できない文字が含まれています",
};
