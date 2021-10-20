import { shallowMount } from "@vue/test-utils";
import Alert from "@/components/Alert.vue";
import { MessageType } from "@/utils/message";

describe("Alert.vue", () => {
  describe("messageプロパティ", () => {
    const props = {
      type: MessageType.Info,
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

  describe("typeプロパティ", () => {
    const props = {
      type: MessageType.Info,
      message: "test",
    };

    describe("MessageTypeの各タイプが指定された場合", () => {
      test("指定されたタイプで表示されること", () => {
        const testType = (t: MessageType) => {
          props.type = t;
          const wrapper = shallowMount(Alert, { propsData: props });
          expect(wrapper.text()).toBe("test");
          expect(wrapper.find("v-alert-stub").attributes("type")).toBe(t);
        };

        testType(MessageType.Info);
        testType(MessageType.Success);
        testType(MessageType.Warn);
        testType(MessageType.Error);
      });
    });
  });
});
