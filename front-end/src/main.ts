import Vue from "vue";
import App from "./App.vue";
import "./registerServiceWorker";
import router from "./router";
import store from "./store";
import vuetify from "./plugins/vuetify";
import axios from "axios";
import { AppMessage, AppMessageType } from "./utils/app-message";

Vue.config.productionTip = false;

// リアクティブなメッセージを設定
const appMsg = Vue.observable<AppMessage>(new AppMessage(AppMessageType.Error, ""));
Vue.prototype.$appMsg = appMsg;

// axiosの設定
axios.defaults.baseURL = process.env.VUE_APP_BACKEND_BASE_URL;
axios.defaults.timeout = 10 * 1000;
Vue.prototype.$http = axios;

// エラーハンドラの設定
Vue.config.errorHandler = (err, vm, info) => {
  if (process.env.NODE_ENV === "development") {
    console.error(err);
    console.error(info);
  }

  if (axios.isAxiosError(err)) {
    if (err.response) {
      // 200番台以外のレスポンスで、メッセージがある場合
      if (err.response.data.message) {
        Vue.prototype.$appMsg.type = AppMessageType.Error;
        Vue.prototype.$appMsg.message = err.response.data.message;
        return;
      }
    } else if (err.request) {
      // サーバからレスポンスを受信できなかった場合
      Vue.prototype.$appMsg.type = AppMessageType.Error;
      Vue.prototype.$appMsg.message = "サーバとの通信に失敗しました";
      return;
    }
  }

  Vue.prototype.$appMsg.type = AppMessageType.Error;
  Vue.prototype.$appMsg.message = "予期せぬエラーが発生しました";
};

new Vue({
  router,
  store,
  vuetify,
  render: (h) => h(App),
}).$mount("#app");
