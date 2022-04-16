import { shallowMount, createLocalVue } from "@vue/test-utils";
import VueRouter from "vue-router";
import App from "@/App.vue";
import AppMessage from "@/components/AppMessage.vue";
import flushPromises from "flush-promises";
import { AppMessageSize } from "@/utils/app-message";
import { paths } from "./test-consts";
import { initAppMsg } from "./test-utils";

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
      const wrapper = shallowMount(App, { localVue });
      expect(wrapper.findComponent(AppMessage).exists()).toBeFalsy();
    });
  });

  describe("メッセージの変更", () => {
    describe("メッセージが表示されていない場合", () => {
      describe("メッセージに空文字以外を設定した場合", () => {
        test("メッセージが表示されること", async () => {
          const wrapper = shallowMount(App, { localVue });

          wrapper.vm.$appMsg.message = "test message";
          await flushPromises();

          expect(wrapper.findComponent(AppMessage).exists()).toBeTruthy();
          expect(wrapper.findComponent(AppMessage).attributes(MESSAGE_PROPS_NAME)).toBe("test message");
        });
      });

      describe("メッセージに空文字を設定した場合", () => {
        test("メッセージが表示されないこと", async () => {
          const wrapper = shallowMount(App, { localVue });

          wrapper.vm.$appMsg.message = "";
          await flushPromises();

          expect(wrapper.findComponent(AppMessage).exists()).toBeFalsy();
        });
      });
    });

    describe("メッセージが表示されている場合", () => {
      describe("メッセージに空文字以外を設定した場合", () => {
        test("メッセージが更新されること", async () => {
          const wrapper = shallowMount(App, { localVue });

          // メッセージ表示
          wrapper.vm.$appMsg.message = "test message";
          await flushPromises();
          expect(wrapper.findComponent(AppMessage).exists()).toBeTruthy();
          expect(wrapper.findComponent(AppMessage).attributes(MESSAGE_PROPS_NAME)).toBe("test message");

          // 表示中のメッセージとは異なるメッセージを設定
          wrapper.vm.$appMsg.message = "updated message";
          await flushPromises();

          expect(wrapper.findComponent(AppMessage).exists()).toBeTruthy();
          expect(wrapper.findComponent(AppMessage).attributes(MESSAGE_PROPS_NAME)).toBe("updated message");
        });
      });

      describe("メッセージに空文字を設定した場合", () => {
        test("メッセージが非表示になること", async () => {
          const wrapper = shallowMount(App, { localVue });

          // メッセージ表示
          wrapper.vm.$appMsg.message = "test message";
          await flushPromises();
          expect(wrapper.findComponent(AppMessage).exists()).toBeTruthy();
          expect(wrapper.findComponent(AppMessage).attributes(MESSAGE_PROPS_NAME)).toBe("test message");

          wrapper.vm.$appMsg.message = "";
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
      const wrapper = shallowMount(App, { localVue });

      wrapper.vm.$appMsg.message = "test message";
      wrapper.vm.$appMsg.setSize(inputSize);
      await flushPromises();

      expect(wrapper.findComponent(AppMessage).exists()).toBeTruthy();
      expect(wrapper.findComponent(AppMessage).attributes(MESSAGE_PROPS_NAME)).toBe("test message");
      expect(wrapper.findComponent({ ref: "appMessageColumn" }).attributes("md")).toBe(mdSize);
      expect(wrapper.findComponent({ ref: "appMessageColumn" }).attributes("lg")).toBe(lgSize);
      expect(wrapper.findComponent({ ref: "appMessageColumn" }).attributes("xl")).toBe(xlSize);
    });
  });

  describe("ページ遷移", () => {
    describe("メッセージが表示されている場合", () => {
      test("メッセージが非表示になること", async () => {
        const router = new VueRouter();
        router.push(paths.Login);
        const wrapper = shallowMount(App, { localVue, router });

        // メッセージ表示
        wrapper.vm.$appMsg.message = "test message";
        await flushPromises();
        expect(wrapper.findComponent(AppMessage).exists()).toBeTruthy();
        expect(wrapper.findComponent(AppMessage).attributes(MESSAGE_PROPS_NAME)).toBe("test message");

        // ページ遷移
        router.push(paths.Join);
        await flushPromises();

        expect(wrapper.findComponent(AppMessage).exists()).toBeFalsy();
      });
    });

    describe("メッセージが表示されていない場合", () => {
      test("メッセージが非表示のままであること", async () => {
        const router = new VueRouter();
        router.push(paths.Login);
        const wrapper = shallowMount(App, { localVue, router });

        // ページ遷移
        router.push(paths.Join);
        await flushPromises();

        expect(wrapper.findComponent(AppMessage).exists()).toBeFalsy();
      });
    });
  });
});
