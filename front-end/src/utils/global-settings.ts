import ja from "vee-validate/dist/locale/ja.json";
import { extend, localize } from "vee-validate";
import { max, min, required } from "vee-validate/dist/rules";
import { customEmail, password } from "./validations";
import axios from "axios";
import Vue from "vue";
import { AppMessageType } from "@/store/app-message";

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

/**
 * エラーハンドラ
 * @param err Errorオブジェクト
 * @param vm Vueオブジェクト
 * @param info エラー情報
 */
export function errorHandler(err: Error, vm: Vue, info: string): void {
  if (process.env.NODE_ENV === "development") {
    console.error(err);
    console.error(info);
  }

  if (axios.isAxiosError(err)) {
    if (err.response) {
      // 200番台以外のレスポンスで、メッセージがある場合
      if (err.response.data.message) {
        Vue.prototype.$state.appMsg.type = AppMessageType.Error;
        Vue.prototype.$state.appMsg.message = err.response.data.message;
        return;
      }
    } else if (err.request) {
      // サーバからレスポンスを受信できなかった場合
      Vue.prototype.$state.appMsg.type = AppMessageType.Error;
      Vue.prototype.$state.appMsg.message = "サーバとの通信に失敗しました";
      return;
    }
  }

  Vue.prototype.$state.appMsg.type = AppMessageType.Error;
  Vue.prototype.$state.appMsg.message = "予期せぬエラーが発生しました";
}
