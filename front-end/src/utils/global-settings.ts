import ja from "vee-validate/dist/locale/ja.json";
import { extend, localize } from "vee-validate";
import { max, min, required } from "vee-validate/dist/rules";
import { customEmail, password } from "./validations";
import axios, { AxiosRequestConfig, AxiosResponse, AxiosStatic } from "axios";
import Vue from "vue";
import { AppMessageType } from "@/store/app-message";
import router from "@/router";
import { StatusCodes } from "http-status-codes";

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
 * axiosの設定を行う
 * @returns axios
 */
export function initAxios(): AxiosStatic {
  axios.defaults.baseURL = process.env.VUE_APP_BACKEND_BASE_URL;
  axios.defaults.timeout = 10 * 1000;
  axios.interceptors.request.use(onRequestFulfilled, onRejected);
  axios.interceptors.response.use(onResponseFulfilled, onResponseRejected);
  return axios;
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

function onRequestFulfilled(config: AxiosRequestConfig): AxiosRequestConfig {
  Vue.prototype.$state.loading = true;
  return config;
}

// eslint-disable-next-line @typescript-eslint/no-explicit-any
function onResponseFulfilled(response: AxiosResponse<any>): AxiosResponse<any> {
  Vue.prototype.$state.loading = false;
  return response;
}

// eslint-disable-next-line @typescript-eslint/no-explicit-any
function onResponseRejected(error: any): Promise<never> {
  // 認証が必要なページで401エラーとなった場合、ログイン画面に遷移
  if (axios.isAxiosError(error) && error.response?.status === StatusCodes.UNAUTHORIZED) {
    router.push("/login");
  }

  return onRejected(error);
}

// eslint-disable-next-line @typescript-eslint/no-explicit-any
function onRejected(error: any): Promise<never> {
  Vue.prototype.$state.loading = false;
  return Promise.reject(error);
}
