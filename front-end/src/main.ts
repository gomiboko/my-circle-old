import Vue from "vue";
import App from "./App.vue";
import "./registerServiceWorker";
import router from "./router";
import vuetify from "./plugins/vuetify";
import axios from "axios";
import { state } from "./store/store";
import { errorHandler, initVeeValidate } from "./utils/global-settings";

Vue.config.productionTip = false;

// 状態管理用のオブジェクトを設定
Vue.prototype.$state = state;

// axiosの設定
axios.defaults.baseURL = process.env.VUE_APP_BACKEND_BASE_URL;
axios.defaults.timeout = 10 * 1000;
axios.interceptors.request.use(
  (config) => {
    Vue.prototype.$state.loading = true;
    return config;
  },
  (error) => {
    Vue.prototype.$state.loading = false;
    return Promise.reject(error);
  }
);
axios.interceptors.response.use(
  (response) => {
    Vue.prototype.$state.loading = false;
    return response;
  },
  (error) => {
    Vue.prototype.$state.loading = false;
    return Promise.reject(error);
  }
);
Vue.prototype.$http = axios;

// エラーハンドラの設定
Vue.config.errorHandler = errorHandler;

// VeeValidateの設定
initVeeValidate();

new Vue({
  router,
  vuetify,
  render: (h) => h(App),
}).$mount("#app");
