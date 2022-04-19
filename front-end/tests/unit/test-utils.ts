import { Wrapper } from "@vue/test-utils";
import { ValidationProvider } from "vee-validate";
import flushPromises from "flush-promises";
import Vue from "vue";
import { state } from "@/store/store";

/**
 * 指定されたコンポーネントの input 要素に値を設定する
 * @param wrapper 入力対象コンポーネントの Wrapper オブジェクト
 * @param val 入力する値
 */
export async function setValue<T extends Vue>(wrapper: Wrapper<T>, val: string): Promise<void> {
  await wrapper.find("input").setValue(val);
  await flushPromises();
}

/**
 * 指定された ref 属性値を持つ validation-provider コンポーネントの errors スロットプロパティを取得する
 * @param wrapper コンポーネントの Wrapper オブジェクト
 * @param providerRefName ref 属性値
 * @returns errors スロットプロパティ
 */
export function getValidationProviderErrors<T extends Vue>(wrapper: Wrapper<T>, providerRefName: string): string[] {
  return (wrapper.vm.$refs[providerRefName] as InstanceType<typeof ValidationProvider>).errors;
}

/**
 * ValidationObserver の状態を更新する為に必要な処理を実行する。
 * 事前に jest.userFakeTimers() を実行しておくこと。
 */
export async function flushAll(): Promise<void> {
  await flushPromises();
  jest.runAllTimers();
  await flushPromises();
}

/**
 * 指定された文字数のメールアドレスを生成する
 * @param length メールアドレス全体の文字数
 * @returns 指定された文字数のメールアドレス
 */
export function createEmailAddress(length: number): string {
  return "a".repeat(length - "@example.com".length) + "@example.com";
}

/**
 * アプリメッセージの初期化を行う
 */
export function initAppMsg(): void {
  Vue.prototype.$state = state;
}
