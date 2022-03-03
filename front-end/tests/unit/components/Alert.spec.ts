import { shallowMount } from "@vue/test-utils";
import Alert from "@/components/Alert.vue";
import { MessageType } from "@/utils/message";

describe("Alert.vue", () => {
  describe("messageプロパティ", () => {
    const props = {
      messageType: MessageType.Info,
      message: "",
    };

    describe("空の場合", () => {
      test("メッセージが表示されないこと", () => {
        props.message = "";
        const wrapper = shallowMount(Alert, { propsData: props });
        expect(wrapper.text()).toBe("");
      });
    });

    describe("空でない場合", () => {
      test("設定されている値が表示されること", () => {
        props.message = "テストメッセージ";
        const wrapper = shallowMount(Alert, { propsData: props });
        expect(wrapper.text()).toBe("テストメッセージ");
      });
    });
  });

  describe("messageTypeプロパティ", () => {
    describe("MessageType列挙型の各タイプが指定された場合", () => {
      test("指定されたタイプで表示されること", () => {
        const props = {
          messageType: MessageType.Info,
          message: "test",
        };

        const testType = (t: MessageType) => {
          props.messageType = t;
          const wrapper = shallowMount(Alert, { propsData: props });
          expect(wrapper.text()).toBe("test");
          expect(wrapper.findComponent({ name: "v-alert" }).attributes("type")).toBe(t);
        };

        testType(MessageType.Info);
        testType(MessageType.Success);
        testType(MessageType.Warn);
        testType(MessageType.Error);
      });
    });
  });
});
