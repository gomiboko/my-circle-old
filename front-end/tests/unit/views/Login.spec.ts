import { shallowMount, mount } from "@vue/test-utils";
import { ValidationObserver, ValidationProvider } from "vee-validate";
import Login from "@/views/Login.vue";
import VueRouter from "vue-router";
import { AxiosError } from "axios";
import { Message, MSG_EVENT } from "@/utils/message";
import { getValidationProviderErrors, flushAll, getEventCount, getVeryFirstEventData, setValue } from "../test-utils";
import flushPromises from "flush-promises";
import { createMockedLocalVue } from "../local-vue";

jest.useFakeTimers();

describe("Login.vue", () => {
  describe("ログインボタンの活性制御", () => {
    test("必須項目の入力状態によって活性制御されること", async () => {
      const wrapper = shallowMount(Login, {
        stubs: { ValidationObserver, ValidationProvider },
      });
      await flushAll();

      const loginBtnWrapper = wrapper.find("v-btn-stub");

      // 初期表示時
      expect(loginBtnWrapper.attributes("disabled")).toBeTruthy();

      // 必須項目の一部が入力された場合
      await wrapper.setData({ email: "a" });
      await flushAll();
      expect(loginBtnWrapper.attributes("disabled")).toBeTruthy();

      // 必須項目が全て入力された場合
      await wrapper.setData({ password: "b" });
      await flushAll();
      expect(loginBtnWrapper.attributes("disabled")).toBeFalsy();

      // 必須項目の一部が削除された場合
      await wrapper.setData({ password: "" });
      await flushAll();
      expect(loginBtnWrapper.attributes("disabled")).toBeTruthy();
    });
  });

  describe("バリデーション", () => {
    describe("メールアドレステキストボックス", () => {
      describe("空の場合", () => {
        test("エラーメッセージが表示されること", async () => {
          const { localVue } = createMockedLocalVue();
          const router = new VueRouter();
          const wrapper = mount(Login, { localVue, router });

          const emailTextWrapper = wrapper.findComponent({ ref: "emailTextField" });
          await setValue(emailTextWrapper, "");
          await flushPromises();

          const errors = getValidationProviderErrors(wrapper, "emailTextFieldProvider");
          expect(errors.length).toBe(1);
        });
      });

      describe("値が入力された場合", () => {
        test("エラーメッセージが表示されないこと", async () => {
          const { localVue } = createMockedLocalVue();
          const router = new VueRouter();
          const wrapper = mount(Login, { localVue, router });

          const emailTextWrapper = wrapper.findComponent({ ref: "emailTextField" });
          await setValue(emailTextWrapper, "a");
          await flushPromises();

          const errors = getValidationProviderErrors(wrapper, "emailTextFieldProvider");
          expect(errors.length).toBe(0);
        });
      });
    });

    describe("パスワードテキストボックス", () => {
      describe("空の場合", () => {
        test("エラーメッセージが表示されること", async () => {
          const { localVue } = createMockedLocalVue();
          const router = new VueRouter();
          const wrapper = mount(Login, { localVue, router });

          const passTextWrapper = wrapper.findComponent({ ref: "passwordTextField" });
          await setValue(passTextWrapper, "");
          await flushPromises();

          const errors = getValidationProviderErrors(wrapper, "passwordTextFieldProvider");
          expect(errors.length).toBe(1);
        });
      });

      describe("値が入力された場合", () => {
        test("エラーメッセージが表示されないこと", async () => {
          const { localVue } = createMockedLocalVue();
          const router = new VueRouter();
          const wrapper = mount(Login, { localVue, router });

          const passTextWrapper = wrapper.findComponent({ ref: "passwordTextField" });
          await setValue(passTextWrapper, "a");
          await flushPromises();

          const errors = getValidationProviderErrors(wrapper, "passwordTextFieldProvider");
          expect(errors.length).toBe(0);
        });
      });
    });
  });

  describe("ログインボタン押下", () => {
    describe("ログインが成功した場合", () => {
      test("トップページに遷移すること", async () => {
        const { localVue, axiosMock } = createMockedLocalVue();

        // ログイン成功時のレスポンスはステータスコードのみ
        axiosMock.post.mockResolvedValue(null);

        const router = new VueRouter();
        router.push("/login");
        const wrapper = mount(Login, { localVue, router });

        const emailTextWrapper = wrapper.findComponent({ ref: "emailTextField" });
        const passTextWrapper = wrapper.findComponent({ ref: "passwordTextField" });
        await setValue(emailTextWrapper, "foo@example.com");
        await setValue(passTextWrapper, "password");
        await flushAll();

        const loginBtnWrapper = wrapper.findComponent({ ref: "loginButton" });
        await loginBtnWrapper.find("button").trigger("click");
        await flushPromises();

        expect(wrapper.vm.$route.path).toBe("/");
      });
    });

    describe("ログインに失敗した場合", () => {
      test("ログインページにエラーメッセージが表示されること", async () => {
        const { localVue, axiosMock } = createMockedLocalVue();

        axiosMock.post.mockRejectedValue({
          isAxiosError: true,
          response: {
            data: { message: "ログイン失敗テスト" },
          },
        } as AxiosError);
        axiosMock.isAxiosError.mockReturnValue(true);

        const router = new VueRouter();
        router.push("/login");
        const w = mount(Login, { localVue, router });

        const emailTextWrapper = w.findComponent({ ref: "emailTextField" });
        const passTextWrapper = w.findComponent({ ref: "passwordTextField" });
        await setValue(emailTextWrapper, "wrong_user");
        await setValue(passTextWrapper, "wrong_password");
        await flushAll();

        const loginBtnWrapper = w.findComponent({ ref: "loginButton" });
        await loginBtnWrapper.find("button").trigger("click");
        await flushPromises();

        // メッセージ表示のカスタムイベントが1回発生していること
        expect(getEventCount(w, MSG_EVENT)).toBe(1);
        // 「予期せぬエラー」のメッセージでないこと
        const eventData = getVeryFirstEventData<Login, Message>(w, MSG_EVENT);
        expect(eventData.message).not.toContain("予期せぬエラー");
        // ページ遷移していないこと
        expect(w.vm.$route.path).toBe("/login");
      });
    });

    describe("予期せぬエラーが発生した場合", () => {
      test("ログインページにエラーメッセージが表示されること", async () => {
        const { localVue, axiosMock } = createMockedLocalVue();

        axiosMock.post.mockRejectedValue(new Error("予期せぬエラーテスト"));
        axiosMock.isAxiosError.mockReturnValue(false);

        const router = new VueRouter();
        router.push("/login");
        const w = mount(Login, { localVue, router });

        const emailTextWrapper = w.findComponent({ ref: "emailTextField" });
        const passTextWrapper = w.findComponent({ ref: "passwordTextField" });
        await setValue(emailTextWrapper, "foo@example.com");
        await setValue(passTextWrapper, "password");
        await flushAll();

        const loginBtnWrapper = w.findComponent({ ref: "loginButton" });
        await loginBtnWrapper.find("button").trigger("click");
        await flushPromises();

        // メッセージ表示のカスタムイベントが1回発生していること
        expect(getEventCount(w, MSG_EVENT)).toBe(1);
        // 「予期せぬエラー」のメッセージであること
        const eventData = getVeryFirstEventData<Login, Message>(w, MSG_EVENT);
        expect(eventData.message).toContain("予期せぬエラー");
        // ページ遷移していないこと
        expect(w.vm.$route.path).toBe("/login");
      });
    });
  });
});
