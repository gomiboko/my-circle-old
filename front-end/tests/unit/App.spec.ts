import { shallowMount, createLocalVue } from "@vue/test-utils";
import VueRouter from "vue-router";
import App from "@/App.vue";
import { Message, MessageType, MSG_EVENT } from "@/utils/message";
import Alert from "@/components/Alert.vue";

// App.vue の message データオブジェクト名
const MESSAGE_DATA_NAME = "message";

// Alert.vue の message プロパティ名
const MESSAGE_PROPS_NAME = "message";

// router-view コンポーネントのスタブ要素名
const ROUTER_VIEW_STUB_NAME = "router-view-stub";

const localVue = createLocalVue();
localVue.use(VueRouter);

describe("App.vue", () => {
  describe("初期表示", () => {
    test("メッセージが表示されていないこと", () => {
      const wrapper = shallowMount(App, { localVue });
      expect(wrapper.findComponent(Alert).exists()).toBeFalsy();
    });
  });

  describe("メッセージ表示イベント", () => {
    describe("メッセージが表示されていない場合", () => {
      describe("イベントのメッセージが空でない場合", () => {
        test("メッセージが表示されること", async () => {
          const wrapper = shallowMount(App, { localVue });

          const msg = new Message(MessageType.Info, "test message");
          await wrapper.find(ROUTER_VIEW_STUB_NAME).trigger(MSG_EVENT, msg);

          expect(
            wrapper.findComponent(Alert).attributes(MESSAGE_PROPS_NAME)
          ).toBe("test message");
          expect(true).toBe(false);
        });
      });

      describe("イベントのメッセージが空の場合", () => {
        test("メッセージが表示されないこと", async () => {
          const wrapper = shallowMount(App, { localVue });

          const msg = new Message(MessageType.Info, "");
          await wrapper.find(ROUTER_VIEW_STUB_NAME).trigger(MSG_EVENT, msg);

          expect(wrapper.findComponent(Alert).exists()).toBeFalsy();
        });
      });
    });

    describe("メッセージが表示されている場合", () => {
      describe("イベントのメッセージが空でない場合", () => {
        test("メッセージが更新されること", async () => {
          const wrapper = shallowMount(App, { localVue });

          // メッセージ表示
          const msg = new Message(MessageType.Info, "test message");
          await wrapper.find(ROUTER_VIEW_STUB_NAME).trigger(MSG_EVENT, msg);
          expect(
            wrapper.findComponent(Alert).attributes(MESSAGE_PROPS_NAME)
          ).toBe("test message");

          // 表示中のメッセージとは異なるメッセージでイベントを発火
          msg.message = "updated message";
          await wrapper.find(ROUTER_VIEW_STUB_NAME).trigger(MSG_EVENT, msg);

          expect(
            wrapper.findComponent(Alert).attributes(MESSAGE_PROPS_NAME)
          ).toBe("updated message");
        });
      });

      describe("イベントのメッセージが空の場合", () => {
        test("メッセージが非表示になること", async () => {
          const wrapper = shallowMount(App, { localVue });

          // メッセージ表示
          const msg = new Message(MessageType.Info, "test message");
          await wrapper.find(ROUTER_VIEW_STUB_NAME).trigger(MSG_EVENT, msg);
          expect(
            wrapper.findComponent(Alert).attributes(MESSAGE_PROPS_NAME)
          ).toBe("test message");

          msg.message = "";
          await wrapper.find(ROUTER_VIEW_STUB_NAME).trigger(MSG_EVENT, msg);

          expect(wrapper.findComponent(Alert).exists()).toBeFalsy();
        });
      });
    });
  });

  describe("ページ遷移", () => {
    describe("メッセージが表示されている場合", () => {
      test("メッセージが非表示になること", async () => {
        const wrapper = shallowMount(App, { localVue });

        // メッセージ表示
        const msg = new Message(MessageType.Info, "test message");
        await wrapper.find(ROUTER_VIEW_STUB_NAME).trigger(MSG_EVENT, msg);
        expect(
          wrapper.findComponent(Alert).attributes(MESSAGE_PROPS_NAME)
        ).toBe("test message");

        // ページ遷移時の処理(messageデータオブジェクトに空文字を設定)を実行
        await wrapper.setData({ [MESSAGE_DATA_NAME]: "" });

        expect(wrapper.findComponent(Alert).exists()).toBeFalsy();
      });
    });

    describe("メッセージが表示されていない場合", () => {
      test("メッセージが非表示のままであること", async () => {
        const wrapper = shallowMount(App, { localVue });

        // ページ遷移時の処理(messageデータオブジェクトに空文字を設定)を実行
        await wrapper.setData({ [MESSAGE_DATA_NAME]: "" });

        expect(wrapper.findComponent(Alert).exists()).toBeFalsy();
      });
    });
  });
});
