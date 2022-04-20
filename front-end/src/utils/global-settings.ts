import ja from "vee-validate/dist/locale/ja.json";
import { extend, localize } from "vee-validate";
import { max, min, required } from "vee-validate/dist/rules";
import { customEmail, password } from "./validations";

/**
 * VeeValidateの設定を行う
 */
export function initVeeValidate(): void {
  localize("ja", ja);
  extend("required", required);
  extend("min", min);
  extend("max", max);
  extend("email", customEmail);
  extend("password", password);
}
