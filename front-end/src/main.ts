import Vue from "vue";
import App from "./App.vue";
import "./registerServiceWorker";
import router from "./router";
import vuetify from "./plugins/vuetify";
import { state } from "./store/store";
import { errorHandler, initVeeValidate, initAxios } from "./utils/global-settings";

Vue.config.productionTip = false;

// 状態管理用のオブジェクトを設定
Vue.prototype.$state = state;

// axiosの設定
Vue.prototype.$http = initAxios();

// エラーハンドラの設定
Vue.config.errorHandler = errorHandler;

// VeeValidateの設定
initVeeValidate();

new Vue({
  router,
  vuetify,
  render: (h) => h(App),
}).$mount("#app");
