import { shallowMount, mount } from "@vue/test-utils";
import { ValidationObserver, ValidationProvider } from "vee-validate";
import Join from "@/views/Join.vue";
import {
  flushAll,
  getValidationProviderErrors,
  setValue,
  createEmailAddress,
  getEventCount,
  getVeryFirstEventData,
} from "../test-utils";
import { consts, lengths, paths } from "../test-consts";
import { createMockedLocalVue } from "../local-vue";
import VueRouter from "vue-router";
import Vuetify from "vuetify";
import flushPromises from "flush-promises";
import { AxiosError } from "axios";
import { Message, MSG_EVENT } from "@/utils/message";

const RefUsernameTextField = "usernameTextField";
const RefPasswordTextField = "passwordTextField";
const RefRegisterButton = "registerButton";
const RefUsernameValidationProvider = "usernameTextFieldProvider";
const RefEmailValidationProvider = "emailTextFieldProvider";
const RefPasswordValidationProvider = "passwordTextFieldProvider";

jest.useFakeTimers();

describe("Join.vue", () => {
  describe("パスワードテキストボックスのアイコン表示", () => {
    describe("初期表示", () => {
      test("パスワード非表示のアイコンとなっていること", async () => {
        const wrapper = shallowMount(Join, {
          stubs: { ValidationObserver, ValidationProvider },
        });
        await flushAll();

        expect(wrapper.findComponent({ ref: RefPasswordTextField }).attributes("append-icon")).toBe("mdi-eye-off");
      });
    });

    describe("アイコンクリック時", () => {
      test("非表示アイコンと表示アイコンが切り替わること", async () => {
        const wrapper = shallowMount(Join, {
          stubs: { ValidationObserver, ValidationProvider },
        });
        await flushAll();

        const passTextField = wrapper.findComponent({ ref: RefPasswordTextField });

        // アイコンをクリックすると表示アイコンになること
        passTextField.vm.$emit("click:append");
        await flushPromises();
        expect(passTextField.attributes("append-icon")).toBe("mdi-eye");

        // 再度アイコンをクリックすると非表示アイコンになること
        passTextField.vm.$emit("click:append");
        await flushPromises();
        expect(passTextField.attributes("append-icon")).toBe("mdi-eye-off");
      });
    });
  });

  describe("バリデーション", () => {
    let vuetify: Vuetify;

    beforeEach(() => (vuetify = new Vuetify()));

    const mountWithNewVuetify = () => mount(Join, { vuetify });

    describe("ユーザ名テキストボックス", () => {
      const UsernameIsRequired = "ユーザ名は必須項目です";

      test.each([
        ["空の場合", "", 1, UsernameIsRequired],
        ["半角スペースが入力された場合", consts.HalfWidthSpace, 1, UsernameIsRequired],
        ["全角スペースが入力された場合", consts.FullWidthSpace, 1, UsernameIsRequired],
        ["1文字入力された場合", "a", 0, undefined],
        ["45文字入力された場合", "a".repeat(lengths.UserNameMaxLength), 0, undefined],
        ["46文字入力された場合", "a".repeat(lengths.UserNameMaxLength + 1), 1, "ユーザ名は45文字以内にしてください"],
      ])("%s", async (explanation, input, errCnt, errMsg) => {
        const wrapper = mountWithNewVuetify();
        const usernameTextWrapper = wrapper.findComponent({ ref: RefUsernameTextField });
        await setValue(usernameTextWrapper, input);

        const errors = getValidationProviderErrors(wrapper, RefUsernameValidationProvider);
        expect(errors.length).toBe(errCnt);
        expect(errors[0]).toBe(errMsg);
      });
    });

    describe("メールアドレステキストボックス", () => {
      const RefEmailTextField = "emailTextField";
      const EmailIsRequired = "メールアドレスは必須項目です";
      const InvalidEmail = "有効なメールアドレスではありません";

      // prettier-ignore
      test.each([
        ["空の場合", "", 1, EmailIsRequired],
        ["半角スペースが入力された場合", consts.HalfWidthSpace, 1, EmailIsRequired],
        ["全角スペースが入力された場合", consts.FullWidthSpace, 1, EmailIsRequired],
        ["メールアドレスではない文字列が入力された場合", "a", 1, InvalidEmail],
        ["254文字のメールアドレスが入力された場合", createEmailAddress(lengths.EmailMaxLength), 0, undefined],
        ["255文字のメールアドレスが入力された場合", createEmailAddress(lengths.EmailMaxLength + 1), 1, "メールアドレスは254文字以内にしてください"],
      ])("%s", async (explanation, input, errCnt, errMsg) => {
        const wrapper = mountWithNewVuetify();
        const emailTextWrapper = wrapper.findComponent({ ref: RefEmailTextField });
        await setValue(emailTextWrapper, input);

        const errors = getValidationProviderErrors(wrapper, RefEmailValidationProvider);
        expect(errors.length).toBe(errCnt);
        expect(errors[0]).toBe(errMsg);
      });
    });

    describe("パスワードテキストボックス", () => {
      const InvalidCharacterContained = "パスワードに使用できない文字が含まれています";

      // prettier-ignore
      test.each([
        ["空の場合", "", 1, "パスワードは必須項目です"],
        ["7文字入力された場合", "a".repeat(lengths.PasswordMinLength - 1), 1, "パスワードは8文字以上でなければなりません"],
        ["8文字入力された場合", "a".repeat(lengths.PasswordMinLength), 0, undefined],
        ["128文字入力された場合", "a".repeat(lengths.PasswordMaxLength), 0, undefined],
        ["129文字入力された場合", "a".repeat(lengths.PasswordMaxLength + 1), 1, "パスワードは128文字以内にしてください"],
        ["半角記号が入力された場合", consts.HalfWidthSymbol, 0, undefined],
        ["半角数字が入力された場合", "1234567890", 0, undefined],
        ["半角英字が入力された場合", "abcdefghijklmnopqrstuvwxyz", 0, undefined],
        ["半角スペースが含まれる場合", `${consts.ValidPassword}${consts.HalfWidthSpace}`, 1, InvalidCharacterContained],
        ["全角スペースが含まれる場合", `${consts.ValidPassword}${consts.FullWidthSpace}`, 1, InvalidCharacterContained],
        ["全角アルファベットが入力された場合", `${consts.ValidPassword}${consts.FullWidthA}`, 1, InvalidCharacterContained],
        ["全角数字が入力された場合", `${consts.ValidPassword}${consts.FullWidth1}`, 1, InvalidCharacterContained],
        ["日本語が入力された場合", `${consts.ValidPassword}あ`, 1, InvalidCharacterContained],
        ["絵文字が入力された場合", `${consts.ValidPassword}😋`, 1, InvalidCharacterContained],
      ])("%s", async (explanation, input, errCnt, errMsg) => {
        const wrapper = mountWithNewVuetify();
        const passTextWrapper = wrapper.findComponent({ ref: RefPasswordTextField });
        await setValue(passTextWrapper, input);

        const errors = getValidationProviderErrors(wrapper, RefPasswordValidationProvider);
        expect(errors.length).toBe(errCnt);
        expect(errors[0]).toBe(errMsg);
      });
    });
  });

  describe("登録ボタン押下", () => {
    describe("入力値エラーがある場合", () => {
      test("エラーメッセージが表示されること", async () => {
        const { localVue, axiosMock } = createMockedLocalVue();

        axiosMock.post.mockResolvedValue(null);

        const router = new VueRouter();
        router.push(paths.Join);

        const wrapper = shallowMount(Join, {
          localVue,
          router,
          stubs: { ValidationObserver, ValidationProvider },
        });

        // 全て未入力で登録ボタン押下
        wrapper.findComponent({ ref: RefRegisterButton }).vm.$emit("click");
        await flushPromises();

        const usernameErrs = getValidationProviderErrors(wrapper, RefUsernameValidationProvider);
        const emailErrs = getValidationProviderErrors(wrapper, RefEmailValidationProvider);
        const passErrs = getValidationProviderErrors(wrapper, RefPasswordValidationProvider);
        expect(wrapper.vm.$route.path).toBe(paths.Join);
        expect(usernameErrs.length).toBe(1);
        expect(emailErrs.length).toBe(1);
        expect(passErrs.length).toBe(1);
      });

      describe("一部の入力値エラーが解消された場合", () => {
        test("エラーメッセージが非表示となること", async () => {
          const { localVue, axiosMock } = createMockedLocalVue();

          axiosMock.post.mockResolvedValue(null);

          const router = new VueRouter();
          router.push(paths.Join);

          const vuetify = new Vuetify();

          // テキストボックスに値を入力する必要があるのでshallowMountではなくmount
          const wrapper = mount(Join, { localVue, router, vuetify });

          // 全て未入力で登録ボタンを押下し、エラーメッセージを表示させる
          const registerBtnWrapper = wrapper.findComponent({ ref: RefRegisterButton });
          registerBtnWrapper.vm.$emit("click");
          await flushPromises();
          let usernameErrs = getValidationProviderErrors(wrapper, RefUsernameValidationProvider);
          let emailErrs = getValidationProviderErrors(wrapper, RefEmailValidationProvider);
          let passErrs = getValidationProviderErrors(wrapper, RefPasswordValidationProvider);
          expect(usernameErrs.length).toBe(1);
          expect(emailErrs.length).toBe(1);
          expect(passErrs.length).toBe(1);

          // ユーザ名に適切な値を入力し、登録ボタンを押下
          const usernameTextWrapper = wrapper.findComponent({ ref: RefUsernameTextField });
          await setValue(usernameTextWrapper, consts.ValidUsername);
          await flushAll();
          registerBtnWrapper.vm.$emit("click");
          await flushPromises();

          usernameErrs = getValidationProviderErrors(wrapper, RefUsernameValidationProvider);
          emailErrs = getValidationProviderErrors(wrapper, RefEmailValidationProvider);
          passErrs = getValidationProviderErrors(wrapper, RefPasswordValidationProvider);
          expect(wrapper.vm.$route.path).toBe(paths.Join);
          expect(usernameErrs.length).toBe(0);
          expect(emailErrs.length).toBe(1);
          expect(passErrs.length).toBe(1);
        });
      });
    });

    describe("登録に成功した場合", () => {
      test("トップページに遷移すること", async () => {
        const { localVue, axiosMock } = createMockedLocalVue();

        axiosMock.post.mockResolvedValue({
          response: {
            data: {
              user: {
                ID: 1,
                Name: consts.ValidUsername,
                Email: consts.ValidEmail,
                PasswordHash: "",
                CreatedAt: "2022-02-27T06:48:47.277Z",
                UpdatedAt: "2022-02-27T06:48:47.277Z",
              },
            },
          },
        });

        const router = new VueRouter();
        router.push(paths.Join);
        const wrapper = shallowMount(Join, {
          localVue,
          router,
          stubs: { ValidationObserver, ValidationProvider },
        });

        await wrapper.setData({ username: consts.ValidUsername });
        await wrapper.setData({ email: consts.ValidEmail });
        await wrapper.setData({ password: consts.ValidPassword });

        wrapper.findComponent({ ref: RefRegisterButton }).vm.$emit("click");
        await flushPromises();

        expect(wrapper.vm.$route.path).toBe(paths.Root);
      });
    });

    describe("登録に失敗した場合", () => {
      test("アカウント作成ページにエラーメッセージが表示されること", async () => {
        const { localVue, axiosMock } = createMockedLocalVue();

        axiosMock.post.mockRejectedValue({
          isAxiosError: true,
          response: {
            data: { message: "登録失敗テスト" },
          },
        } as AxiosError);
        axiosMock.isAxiosError.mockReturnValue(true);

        const router = new VueRouter();
        router.push(paths.Join);
        const wrapper = shallowMount(Join, {
          localVue,
          router,
          stubs: { ValidationObserver, ValidationProvider },
        });

        await wrapper.setData({ username: consts.ValidUsername });
        await wrapper.setData({ email: consts.ValidEmail });
        await wrapper.setData({ password: consts.ValidPassword });
        await flushPromises();

        wrapper.findComponent({ ref: RefRegisterButton }).vm.$emit("click");
        await flushPromises();

        // メッセージ表示のカスタムイベントが1回発生していること
        expect(getEventCount(wrapper, MSG_EVENT)).toBe(1);
        // 「予期せぬエラー」のメッセージでないこと
        const eventData = getVeryFirstEventData<Join, Message>(wrapper, MSG_EVENT);
        expect(eventData.message).not.toContain("予期せぬエラー");
        // ページ遷移していないこと
        expect(wrapper.vm.$route.path).toBe(paths.Join);
      });
    });

    describe("予期せぬエラーが発生した場合", () => {
      test("アカウント作成ページにエラーメッセージが表示されること", async () => {
        const { localVue, axiosMock } = createMockedLocalVue();

        axiosMock.post.mockRejectedValue(new Error("エラーテスト"));
        axiosMock.isAxiosError.mockReturnValue(false);

        const router = new VueRouter();
        router.push(paths.Join);
        const wrapper = shallowMount(Join, {
          localVue,
          router,
          stubs: { ValidationObserver, ValidationProvider },
        });

        await wrapper.setData({ username: consts.ValidUsername });
        await wrapper.setData({ email: consts.ValidEmail });
        await wrapper.setData({ password: consts.ValidPassword });
        await flushPromises();

        wrapper.findComponent({ ref: RefRegisterButton }).vm.$emit("click");
        await flushPromises();

        // メッセージ表示のカスタムイベントが1回発生していること
        expect(getEventCount(wrapper, MSG_EVENT)).toBe(1);
        // 「予期せぬエラー」のメッセージであること
        const eventData = getVeryFirstEventData<Join, Message>(wrapper, MSG_EVENT);
        expect(eventData.message).toContain("予期せぬエラー");
        // ページ遷移していないこと
        expect(wrapper.vm.$route.path).toBe(paths.Join);
      });
    });
  });
});
