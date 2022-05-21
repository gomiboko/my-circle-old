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
];

const router = new VueRouter({
  mode: "history",
  base: process.env.BASE_URL,
  routes,
});

// ナビゲーションガード
router.beforeEach((to, from, next) => {
  Vue.prototype.$state.appMsg.message = "";
  next();
});

export default router;
