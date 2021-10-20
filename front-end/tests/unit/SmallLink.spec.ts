import { shallowMount, createLocalVue } from "@vue/test-utils";
import VueRouter from "vue-router";
import SmallLink from "@/components/SmallLink.vue";

const localVue = createLocalVue();
localVue.use(VueRouter);

describe("SmallLink.vue", () => {
  describe("textプロパティ", () => {
    const props = {
      to: "/",
      text: "",
    };

    describe("空の場合", () => {
      test("何も表示されないこと", () => {
        props.text = "";
        const wrapper = shallowMount(SmallLink, {
          localVue,
          propsData: props,
        });
        expect(wrapper.text()).toBe("");
      });
    });

    describe("空でない場合", () => {
      test("textプロパティの値が表示されること", () => {
        props.text = "テストテキスト";
        const wrapper = shallowMount(SmallLink, {
          localVue,
          propsData: props,
        });
        expect(wrapper.text()).toBe("テストテキスト");
      });
    });
  });

  describe("toプロパティ", () => {
    const props = {
      to: "",
      text: "テストリンク",
    };

    describe("空の場合", () => {
      test("リンク先に何も設定されていないこと", () => {
        props.to = "";
        const wrapper = shallowMount(SmallLink, {
          localVue,
          propsData: props,
        });
        expect(wrapper.find("router-link-stub").attributes("to")).toBe("");
      });
    });

    describe("空でない場合", () => {
      test("リンク先がtoプロパティの値になっていること", () => {
        props.to = "/test/path";
        const wrapper = shallowMount(SmallLink, {
          localVue,
          propsData: props,
        });
        expect(wrapper.find("router-link-stub").attributes("to")).toBe(
          "/test/path"
        );
      });
    });
  });
});
