import { shallowMount } from "@vue/test-utils";
import AppMessage from "@/components/AppMessage.vue";
import { AppMessageType } from "@/utils/app-message";

describe("AppMessage.vue", () => {
  describe("messageプロパティ", () => {
    const props = {
      messageType: AppMessageType.Info,
      message: "",
    };

    describe("空の場合", () => {
      test("メッセージが表示されないこと", () => {
        props.message = "";
        const wrapper = shallowMount(AppMessage, { propsData: props });
        expect(wrapper.text()).toBe("");
      });
    });

    describe("空でない場合", () => {
      test("設定されている値が表示されること", () => {
        props.message = "テストメッセージ";
        const wrapper = shallowMount(AppMessage, { propsData: props });
        expect(wrapper.text()).toBe("テストメッセージ");
      });
    });
  });

  describe("messageTypeプロパティ", () => {
    describe("MessageType列挙型の各タイプが指定された場合", () => {
      test("指定されたタイプで表示されること", () => {
        const props = {
          messageType: AppMessageType.Info,
          message: "test",
        };

        const testType = (t: AppMessageType) => {
          props.messageType = t;
          const wrapper = shallowMount(AppMessage, { propsData: props });
          expect(wrapper.text()).toBe("test");
          expect(wrapper.findComponent({ name: "v-alert" }).attributes("type")).toBe(t);
        };

        testType(AppMessageType.Info);
        testType(AppMessageType.Success);
        testType(AppMessageType.Warn);
        testType(AppMessageType.Error);
      });
    });
  });
});
