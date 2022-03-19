import { shallowMount, mount } from "@vue/test-utils";
import { ValidationObserver, ValidationProvider } from "vee-validate";
import Login from "@/views/Login.vue";
import VueRouter from "vue-router";
import { AxiosError } from "axios";
import { Message, MSG_EVENT } from "@/utils/message";
import { getValidationProviderErrors, flushAll, getEventCount, getVeryFirstEventData, setValue } from "../test-utils";
import flushPromises from "flush-promises";
import { createMockedLocalVue } from "../local-vue";
import { consts, paths } from "../test-consts";

const RefEmailTextField = "emailTextField";
const RefPasswordTextField = "passwordTextField";
const RefLoginButton = "loginButton";
const RefEmailTextFieldProvider = "emailTextFieldProvider";
const RefPasswordTextFieldProvider = "passwordTextFieldProvider";

jest.useFakeTimers();

describe("Login.vue", () => {
  describe("バリデーション", () => {
    describe("メールアドレステキストボックス", () => {
      describe("空の場合", () => {
        test("エラーメッセージが表示されること", async () => {
          const { localVue } = createMockedLocalVue();
          const router = new VueRouter();
          const wrapper = mount(Login, { localVue, router });

          const emailTextWrapper = wrapper.findComponent({ ref: RefEmailTextField });
          await setValue(emailTextWrapper, "");
          await flushPromises();

          const errors = getValidationProviderErrors(wrapper, RefEmailTextFieldProvider);
          expect(errors.length).toBe(1);
        });
      });

      describe("値が入力された場合", () => {
        test("エラーメッセージが表示されないこと", async () => {
          const { localVue } = createMockedLocalVue();
          const router = new VueRouter();
          const wrapper = mount(Login, { localVue, router });

          const emailTextWrapper = wrapper.findComponent({ ref: RefEmailTextField });
          await setValue(emailTextWrapper, "a");
          await flushPromises();

          const errors = getValidationProviderErrors(wrapper, RefEmailTextFieldProvider);
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

          const passTextWrapper = wrapper.findComponent({ ref: RefPasswordTextField });
          await setValue(passTextWrapper, "");
          await flushPromises();

          const errors = getValidationProviderErrors(wrapper, RefPasswordTextFieldProvider);
          expect(errors.length).toBe(1);
        });
      });

      describe("値が入力された場合", () => {
        test("エラーメッセージが表示されないこと", async () => {
          const { localVue } = createMockedLocalVue();
          const router = new VueRouter();
          const wrapper = mount(Login, { localVue, router });

          const passTextWrapper = wrapper.findComponent({ ref: RefPasswordTextField });
          await setValue(passTextWrapper, "a");
          await flushPromises();

          const errors = getValidationProviderErrors(wrapper, RefPasswordTextFieldProvider);
          expect(errors.length).toBe(0);
        });
      });
    });
  });

  describe("ログインボタン押下", () => {
    describe("入力値エラーがある場合", () => {
      test("エラーメッセージが表示されること", async () => {
        const { localVue, axiosMock } = createMockedLocalVue();

        axiosMock.post.mockResolvedValue(null);

        const router = new VueRouter();
        router.push(paths.Login);
        const wrapper = shallowMount(Login, {
          localVue,
          router,
          stubs: { ValidationObserver, ValidationProvider },
        });

        // 全て未入力でログインボタン押下
        const loginBtnWrapper = wrapper.findComponent({ ref: RefLoginButton });
        loginBtnWrapper.vm.$emit("click");
        await flushPromises();

        const emailErrs = getValidationProviderErrors(wrapper, RefEmailTextFieldProvider);
        const passErrs = getValidationProviderErrors(wrapper, RefPasswordTextFieldProvider);
        expect(wrapper.vm.$route.path).toBe(paths.Login);
        expect(emailErrs.length).toBe(1);
        expect(passErrs.length).toBe(1);
      });
    });

    describe("一部の入力値エラーが解消された場合", () => {
      test("エラーメッセージが非表示となること", async () => {
        const { localVue, axiosMock } = createMockedLocalVue();

        axiosMock.post.mockResolvedValue(null);

        const router = new VueRouter();
        router.push(paths.Login);
        const wrapper = mount(Login, { localVue, router });

        // 全て未入力でログインボタンを押下し、エラーメッセージを表示させる
        const loginBtnWrapper = wrapper.findComponent({ ref: RefLoginButton });
        loginBtnWrapper.vm.$emit("click");
        await flushPromises();
        let emailErrs = getValidationProviderErrors(wrapper, RefEmailTextFieldProvider);
        let passErrs = getValidationProviderErrors(wrapper, RefPasswordTextFieldProvider);
        expect(emailErrs.length).toBe(1);
        expect(passErrs.length).toBe(1);

        // メールアドレスに適切な値を入力し、ログインボタンを押下
        const emailTextWrapper = wrapper.findComponent({ ref: RefEmailTextField });
        await setValue(emailTextWrapper, consts.ValidEmail);
        await flushAll();
        loginBtnWrapper.vm.$emit("click");
        await flushPromises();

        emailErrs = getValidationProviderErrors(wrapper, RefEmailTextFieldProvider);
        passErrs = getValidationProviderErrors(wrapper, RefPasswordTextFieldProvider);
        expect(wrapper.vm.$route.path).toBe(paths.Login);
        expect(emailErrs.length).toBe(0);
        expect(passErrs.length).toBe(1);
      });
    });

    describe("ログインが成功した場合", () => {
      test("トップページに遷移すること", async () => {
        const { localVue, axiosMock } = createMockedLocalVue();

        // ログイン成功時のレスポンスはステータスコードのみ
        axiosMock.post.mockResolvedValue(null);

        const router = new VueRouter();
        router.push(paths.Login);
        const wrapper = mount(Login, { localVue, router });

        const emailTextWrapper = wrapper.findComponent({ ref: RefEmailTextField });
        const passTextWrapper = wrapper.findComponent({ ref: RefPasswordTextField });
        await setValue(emailTextWrapper, consts.ValidEmail);
        await setValue(passTextWrapper, consts.ValidPassword);
        await flushAll();

        const loginBtnWrapper = wrapper.findComponent({ ref: RefLoginButton });
        loginBtnWrapper.vm.$emit("click");
        await flushPromises();

        expect(wrapper.vm.$route.path).toBe(paths.Root);
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
        router.push(paths.Login);
        const wrapper = mount(Login, { localVue, router });

        const emailTextWrapper = wrapper.findComponent({ ref: RefEmailTextField });
        const passTextWrapper = wrapper.findComponent({ ref: RefPasswordTextField });
        await setValue(emailTextWrapper, "wrong_user");
        await setValue(passTextWrapper, "wrong_password");
        await flushAll();

        const loginBtnWrapper = wrapper.findComponent({ ref: RefLoginButton });
        loginBtnWrapper.vm.$emit("click");
        await flushPromises();

        // メッセージ表示のカスタムイベントが1回発生していること
        expect(getEventCount(wrapper, MSG_EVENT)).toBe(1);
        // 「予期せぬエラー」のメッセージでないこと
        const eventData = getVeryFirstEventData<Login, Message>(wrapper, MSG_EVENT);
        expect(eventData.message).not.toContain("予期せぬエラー");
        // ページ遷移していないこと
        expect(wrapper.vm.$route.path).toBe(paths.Login);
      });
    });

    describe("予期せぬエラーが発生した場合", () => {
      test("ログインページにエラーメッセージが表示されること", async () => {
        const { localVue, axiosMock } = createMockedLocalVue();

        axiosMock.post.mockRejectedValue(new Error("エラーテスト"));
        axiosMock.isAxiosError.mockReturnValue(false);

        const router = new VueRouter();
        router.push(paths.Login);
        const wrapper = mount(Login, { localVue, router });

        const emailTextWrapper = wrapper.findComponent({ ref: RefEmailTextField });
        const passTextWrapper = wrapper.findComponent({ ref: RefPasswordTextField });
        await setValue(emailTextWrapper, consts.ValidEmail);
        await setValue(passTextWrapper, consts.ValidPassword);
        await flushAll();

        const loginBtnWrapper = wrapper.findComponent({ ref: RefLoginButton });
        loginBtnWrapper.vm.$emit("click");
        await flushPromises();

        // メッセージ表示のカスタムイベントが1回発生していること
        expect(getEventCount(wrapper, MSG_EVENT)).toBe(1);
        // 「予期せぬエラー」のメッセージであること
        const eventData = getVeryFirstEventData<Login, Message>(wrapper, MSG_EVENT);
        expect(eventData.message).toContain("予期せぬエラー");
        // ページ遷移していないこと
        expect(wrapper.vm.$route.path).toBe(paths.Login);
      });
    });
  });

  describe("新規アカウント登録ボタン押下", () => {
    test("ユーザ登録画面が表示されること", async () => {
      const { localVue } = createMockedLocalVue();

      const router = new VueRouter();
      router.push(paths.Login);

      const wrapper = shallowMount(Login, { localVue, router });

      const regAccountBtnWrapper = wrapper.findComponent({ ref: "registerAccountButton" });
      regAccountBtnWrapper.vm.$emit("click");

      expect(wrapper.vm.$route.path).toBe(paths.Join);
    });
  });
});
