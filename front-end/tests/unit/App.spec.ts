import { shallowMount, createLocalVue } from "@vue/test-utils";
import VueRouter from "vue-router";
import App from "@/App.vue";
import { Message, MessageType } from "@/utils/message";
import AppMessage from "@/components/AppMessage.vue";
import { AppMsgSize } from "@/utils/consts";
import flushPromises from "flush-promises";
import { execMethod } from "./test-utils";

// App.vue の message データオブジェクト名
const MESSAGE_DATA_NAME = "message";

// AppMessage.vue の message プロパティ名
const MESSAGE_PROPS_NAME = "message";

// メッセージ表示イベントハンドラ名
const MSG_EVENT_HANDLER_NAME = "showMessage";

const localVue = createLocalVue();
localVue.use(VueRouter);

describe("App.vue", () => {
  describe("初期表示", () => {
    test("メッセージが表示されていないこと", () => {
      const wrapper = shallowMount(App, { localVue });
      expect(wrapper.findComponent(AppMessage).exists()).toBeFalsy();
    });
  });

  describe("メッセージ表示イベント", () => {
    describe("メッセージが表示されていない場合", () => {
      describe("イベントのメッセージが空でない場合", () => {
        test("メッセージが表示されること", async () => {
          const wrapper = shallowMount(App, { localVue });

          const msg = new Message(MessageType.Info, "test message");
          execMethod(wrapper, MSG_EVENT_HANDLER_NAME, msg);
          await flushPromises();

          expect(wrapper.findComponent(AppMessage).exists()).toBeTruthy();
          expect(wrapper.findComponent(AppMessage).attributes(MESSAGE_PROPS_NAME)).toBe("test message");
        });
      });

      describe("イベントのメッセージが空の場合", () => {
        test("メッセージが表示されないこと", async () => {
          const wrapper = shallowMount(App, { localVue });

          const msg = new Message(MessageType.Info, "");
          execMethod(wrapper, MSG_EVENT_HANDLER_NAME, msg);
          await flushPromises();

          expect(wrapper.findComponent(AppMessage).exists()).toBeFalsy();
        });
      });
    });

    describe("メッセージが表示されている場合", () => {
      describe("イベントのメッセージが空でない場合", () => {
        test("メッセージが更新されること", async () => {
          const wrapper = shallowMount(App, { localVue });

          // メッセージ表示
          const msg = new Message(MessageType.Info, "test message");
          execMethod(wrapper, MSG_EVENT_HANDLER_NAME, msg);
          await flushPromises();
          expect(wrapper.findComponent(AppMessage).exists()).toBeTruthy();
          expect(wrapper.findComponent(AppMessage).attributes(MESSAGE_PROPS_NAME)).toBe("test message");

          // 表示中のメッセージとは異なるメッセージでイベントを発火
          msg.message = "updated message";
          execMethod(wrapper, MSG_EVENT_HANDLER_NAME, msg);
          await flushPromises();

          expect(wrapper.findComponent(AppMessage).exists()).toBeTruthy();
          expect(wrapper.findComponent(AppMessage).attributes(MESSAGE_PROPS_NAME)).toBe("updated message");
        });
      });

      describe("イベントのメッセージが空の場合", () => {
        test("メッセージが非表示になること", async () => {
          const wrapper = shallowMount(App, { localVue });

          // メッセージ表示
          const msg = new Message(MessageType.Info, "test message");
          execMethod(wrapper, MSG_EVENT_HANDLER_NAME, msg);
          await flushPromises();
          expect(wrapper.findComponent(AppMessage).exists()).toBeTruthy();
          expect(wrapper.findComponent(AppMessage).attributes(MESSAGE_PROPS_NAME)).toBe("test message");

          msg.message = "";
          execMethod(wrapper, MSG_EVENT_HANDLER_NAME, msg);
          await flushPromises();

          expect(wrapper.findComponent(AppMessage).exists()).toBeFalsy();
        });
      });
    });
  });

  describe("メッセージ表示領域の大きさ", () => {
    test.each([
      ["指定がない場合", undefined, "6"],
      ["デフォルトと同じ大きさが指定された場合", AppMsgSize.Col6, "6"],
      ["デフォルトと異なる大きさが指定された場合", AppMsgSize.Col4, "4"],
    ])("%s", async (explanation, inputSize, expectedSize) => {
      const wrapper = shallowMount(App, { localVue });

      const msg = new Message(MessageType.Info, "test message");
      execMethod(wrapper, MSG_EVENT_HANDLER_NAME, msg, inputSize);
      await flushPromises();

      expect(wrapper.findComponent(AppMessage).exists()).toBeTruthy();
      expect(wrapper.findComponent(AppMessage).attributes(MESSAGE_PROPS_NAME)).toBe("test message");
      expect(wrapper.findComponent({ ref: "appMessageColumn" }).attributes("md")).toBe(expectedSize);
    });
  });

  describe("ページ遷移", () => {
    describe("メッセージが表示されている場合", () => {
      test("メッセージが非表示になること", async () => {
        const wrapper = shallowMount(App, { localVue });

        // メッセージ表示
        const msg = new Message(MessageType.Info, "test message");
        execMethod(wrapper, MSG_EVENT_HANDLER_NAME, msg);
        await flushPromises();
        expect(wrapper.findComponent(AppMessage).exists()).toBeTruthy();
        expect(wrapper.findComponent(AppMessage).attributes(MESSAGE_PROPS_NAME)).toBe("test message");

        // ページ遷移時の処理(messageデータオブジェクトに空文字を設定)を実行
        await wrapper.setData({ [MESSAGE_DATA_NAME]: "" });

        expect(wrapper.findComponent(AppMessage).exists()).toBeFalsy();
      });
    });

    describe("メッセージが表示されていない場合", () => {
      test("メッセージが非表示のままであること", async () => {
        const wrapper = shallowMount(App, { localVue });

        // ページ遷移時の処理(messageデータオブジェクトに空文字を設定)を実行
        await wrapper.setData({ [MESSAGE_DATA_NAME]: "" });

        expect(wrapper.findComponent(AppMessage).exists()).toBeFalsy();
      });
    });
  });
});
