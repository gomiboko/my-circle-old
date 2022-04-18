import Vue from "vue";
import App from "./App.vue";
import "./registerServiceWorker";
import router from "./router";
import vuetify from "./plugins/vuetify";
import axios from "axios";
import { AppMessage, AppMessageType } from "./utils/app-message";
import { errorHandler } from "./utils/global-error-handler";

Vue.config.productionTip = false;

// リアクティブなメッセージを設定
Vue.prototype.$appMsg = Vue.observable<AppMessage>(new AppMessage(AppMessageType.Error, ""));

// axiosの設定
axios.defaults.baseURL = process.env.VUE_APP_BACKEND_BASE_URL;
axios.defaults.timeout = 10 * 1000;
Vue.prototype.$http = axios;

// エラーハンドラの設定
Vue.config.errorHandler = errorHandler;

new Vue({
  router,
  vuetify,
  render: (h) => h(App),
}).$mount("#app");
