import { shallowMount } from "@vue/test-utils";
import RequiredTextField from "@/components/RequiredTextField.vue"

describe("RequiredTextField.vue", () => {
  describe("labelプロパティ", () => {
    test("ラベルに必須マークが付いて表示されること", () => {
      const labelText = "テストラベル";
      const wrapper = shallowMount(RequiredTextField, { propsData: { label: labelText } });
      expect(wrapper.text()).toBe(`${labelText} *`);
    });
  });
});
