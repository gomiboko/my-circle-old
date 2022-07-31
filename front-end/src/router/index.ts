import { PAGE_PATHS } from "@/utils/consts";
import Vue from "vue";
import VueRouter, { RouteConfig } from "vue-router";
import Home from "../views/Home.vue";

Vue.use(VueRouter);

const routes: Array<RouteConfig> = [
  {
    path: PAGE_PATHS.HOME,
    name: "Home",
    component: Home,
  },
  {
    path: "/about",
    name: "About",
    // route level code-splitting
    // this generates a separate chunk (about.[hash].js) for this route
    // which is lazy-loaded when the route is visited.
    component: () => import(/* webpackChunkName: "about" */ "../views/About.vue"),
  },
  {
    path: PAGE_PATHS.LOGIN,
    name: "Login",
    component: () => import("../views/Login.vue"),
  },
  {
    path: PAGE_PATHS.JOIN,
    name: "Join",
    component: () => import("../views/Join.vue"),
  },
  {
    path: PAGE_PATHS.CIRCLE_REGISTER,
    name: "CircleRegister",
    component: () => import("../views/CircleRegister.vue"),
  },
];

const router = new VueRouter({
  mode: "history",
  base: process.env.BASE_URL,
  routes,
});

// ナビゲーションガード
router.beforeEach((to, from, next) => {
  // 遷移後の画面に表示するメッセージがある場合
  if (Vue.prototype.$state.nextScreenMsg.message) {
    Vue.prototype.$state.appMsg.type = Vue.prototype.$state.nextScreenMsg.type;
    Vue.prototype.$state.appMsg.message = Vue.prototype.$state.nextScreenMsg.message;
    Vue.prototype.$state.nextScreenMsg.message = "";
  } else {
    Vue.prototype.$state.appMsg.message = "";
  }

  next();
});

export default router;
