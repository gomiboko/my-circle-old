import { AppMessageType } from "@/store/app-message";
import axios from "axios";
import Vue from "vue";

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
