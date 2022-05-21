import { shallowMount, createLocalVue } from "@vue/test-utils";
import VueRouter, { NavigationGuardNext, Route } from "vue-router";
import App from "@/App.vue";
import AppMessage from "@/components/AppMessage.vue";
import flushPromises from "flush-promises";
import { AppMessageSize } from "@/store/app-message";
import { execAsyncMethod, initAppMsg } from "./test-utils";
import Vue from "vue";
import { createMockedLocalVue } from "./local-vue";
import { API_PATHS, PAGE_PATHS } from "@/utils/consts";

// AppMessage.vue の message プロパティ名
const MESSAGE_PROPS_NAME = "message";

const localVue = createLocalVue();
localVue.use(VueRouter);

beforeEach(() => {
  initAppMsg();
});

describe("App.vue", () => {
  describe("初期表示", () => {
    test("メッセージが表示されていないこと", () => {
      const router = new VueRouter();
      router.push(PAGE_PATHS.HOME);
      const wrapper = shallowMount(App, { localVue, router });

      expect(wrapper.findComponent(AppMessage).exists()).toBeFalsy();
    });
  });

  describe("メッセージの変更", () => {
    describe("メッセージが表示されていない場合", () => {
      describe("メッセージに空文字以外を設定した場合", () => {
        test("メッセージが表示されること", async () => {
          const router = new VueRouter();
          router.push(PAGE_PATHS.HOME);
          const wrapper = shallowMount(App, { localVue, router });

          wrapper.vm.$state.appMsg.message = "test message";
          await flushPromises();

          expect(wrapper.findComponent(AppMessage).exists()).toBeTruthy();
          expect(wrapper.findComponent(AppMessage).attributes(MESSAGE_PROPS_NAME)).toBe("test message");
        });
      });

      describe("メッセージに空文字を設定した場合", () => {
        test("メッセージが表示されないこと", async () => {
          const router = new VueRouter();
          router.push(PAGE_PATHS.HOME);
          const wrapper = shallowMount(App, { localVue, router });

          wrapper.vm.$state.appMsg.message = "";
          await flushPromises();

          expect(wrapper.findComponent(AppMessage).exists()).toBeFalsy();
        });
      });
    });

    describe("メッセージが表示されている場合", () => {
      describe("メッセージに空文字以外を設定した場合", () => {
        test("メッセージが更新されること", async () => {
          const router = new VueRouter();
          router.push(PAGE_PATHS.HOME);
          const wrapper = shallowMount(App, { localVue, router });

          // メッセージ表示
          wrapper.vm.$state.appMsg.message = "test message";
          await flushPromises();
          expect(wrapper.findComponent(AppMessage).exists()).toBeTruthy();
          expect(wrapper.findComponent(AppMessage).attributes(MESSAGE_PROPS_NAME)).toBe("test message");

          // 表示中のメッセージとは異なるメッセージを設定
          wrapper.vm.$state.appMsg.message = "updated message";
          await flushPromises();

          expect(wrapper.findComponent(AppMessage).exists()).toBeTruthy();
          expect(wrapper.findComponent(AppMessage).attributes(MESSAGE_PROPS_NAME)).toBe("updated message");
        });
      });

      describe("メッセージに空文字を設定した場合", () => {
        test("メッセージが非表示になること", async () => {
          const router = new VueRouter();
          router.push(PAGE_PATHS.HOME);
          const wrapper = shallowMount(App, { localVue, router });

          // メッセージ表示
          wrapper.vm.$state.appMsg.message = "test message";
          await flushPromises();
          expect(wrapper.findComponent(AppMessage).exists()).toBeTruthy();
          expect(wrapper.findComponent(AppMessage).attributes(MESSAGE_PROPS_NAME)).toBe("test message");

          wrapper.vm.$state.appMsg.message = "";
          await flushPromises();

          expect(wrapper.findComponent(AppMessage).exists()).toBeFalsy();
        });
      });
    });
  });

  describe("メッセージ表示領域の大きさ", () => {
    test.each([
      ["表示中と同じ大きさが指定された場合", AppMessageSize.Medium, "8", "6", "4"],
      ["表示中と異なる大きさが指定された場合", AppMessageSize.Large, "12", "9", "6"],
    ])("%s", async (explanation, inputSize, mdSize, lgSize, xlSize) => {
      const router = new VueRouter();
      router.push(PAGE_PATHS.HOME);
      const wrapper = shallowMount(App, { localVue, router });

      wrapper.vm.$state.appMsg.message = "test message";
      wrapper.vm.$state.appMsg.setSize(inputSize);
      await flushPromises();

      expect(wrapper.findComponent(AppMessage).exists()).toBeTruthy();
      expect(wrapper.findComponent(AppMessage).attributes(MESSAGE_PROPS_NAME)).toBe("test message");
      expect(wrapper.findComponent({ ref: "appMessageColumn" }).attributes("md")).toBe(mdSize);
      expect(wrapper.findComponent({ ref: "appMessageColumn" }).attributes("lg")).toBe(lgSize);
      expect(wrapper.findComponent({ ref: "appMessageColumn" }).attributes("xl")).toBe(xlSize);
    });
  });

  describe("プロフィールアイコンの表示", () => {
    describe("プロフィールアイコン非表示画面の場合", () => {
      test.each([
        ["ログイン画面の場合", PAGE_PATHS.LOGIN],
        ["アカウント登録画面の場合", PAGE_PATHS.JOIN],
      ])("%s", (explanation, path) => {
        const router = new VueRouter();
        router.push(path);
        const wrapper = shallowMount(App, { localVue, router });

        expect(wrapper.findComponent({ ref: "accountMenu" }).exists()).toBeFalsy();
      });
    });

    describe("プロフィールアイコン表示画面の場合", () => {
      test("プロフィールアイコンが表示されること", () => {
        const router = new VueRouter();
        router.push(PAGE_PATHS.HOME);
        const wrapper = shallowMount(App, { localVue, router });

        expect(wrapper.findComponent({ ref: "accountMenu" }).exists()).toBeTruthy();
      });
    });
  });

  describe("ログアウト処理", () => {
    const LogoutMenuId = 3;

    test("ログアウトAPIが呼ばれ、ログイン画面に遷移すること", async () => {
      const { localVue, axiosMock } = createMockedLocalVue();
      axiosMock.get.mockResolvedValue(null);

      const router = new VueRouter();
      router.push(PAGE_PATHS.HOME);
      const wrapper = shallowMount(App, { localVue, router });

      // ログアウト処理実行
      await execAsyncMethod(wrapper, "onMenuClick", LogoutMenuId);

      expect(axiosMock.get).toBeCalledWith(API_PATHS.LOGOUT, { withCredentials: true });
      expect(wrapper.vm.$route.path).toBe(PAGE_PATHS.LOGIN);
    });
  });

  describe("ページ遷移", () => {
    describe("メッセージが表示されている場合", () => {
      test("メッセージが非表示になること", async () => {
        const router = new VueRouter();
        router.push(PAGE_PATHS.LOGIN);
        router.beforeEach(beforeEachGuard);
        const wrapper = shallowMount(App, { localVue, router });

        // メッセージ表示
        wrapper.vm.$state.appMsg.message = "test message";
        await flushPromises();
        expect(wrapper.findComponent(AppMessage).exists()).toBeTruthy();
        expect(wrapper.findComponent(AppMessage).attributes(MESSAGE_PROPS_NAME)).toBe("test message");

        // ページ遷移
        router.push(PAGE_PATHS.JOIN);
        await flushPromises();

        expect(wrapper.findComponent(AppMessage).exists()).toBeFalsy();
      });
    });

    describe("メッセージが表示されていない場合", () => {
      test("メッセージが非表示のままであること", async () => {
        const router = new VueRouter();
        router.push(PAGE_PATHS.LOGIN);
        router.beforeEach(beforeEachGuard);
        const wrapper = shallowMount(App, { localVue, router });

        // ページ遷移
        router.push(PAGE_PATHS.JOIN);
        await flushPromises();

        expect(wrapper.findComponent(AppMessage).exists()).toBeFalsy();
      });
    });
  });
});

/**
 * Vue Router のナビゲーションガード関数
 * @param to 遷移先のルート情報
 * @param from 繊維元のルート情報
 * @param next ナビゲーションの為のコールバック関数
 */
function beforeEachGuard(to: Route, from: Route, next: NavigationGuardNext<Vue>) {
  Vue.prototype.$state.appMsg.message = "";
  next();
}
