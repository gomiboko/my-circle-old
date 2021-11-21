import { Wrapper } from "@vue/test-utils";
import { ValidationProvider } from "vee-validate";
import flushPromises from "flush-promises";

/**
 * 指定された ref 属性値を持つ validation-provider コンポーネントの errors スロットプロパティを取得する
 * @param wrapper コンポーネントの Wrapper オブジェクト
 * @param providerRefName ref 属性値
 * @returns errors スロットプロパティ
 */
export function getValidationProviderErrors<T extends Vue>(
  wrapper: Wrapper<T>,
  providerRefName: string
): string[] {
  return (
    wrapper.vm.$refs[providerRefName] as InstanceType<typeof ValidationProvider>
  ).errors;
}

/**
 * ValidationObserver の状態を更新する為に必要な処理を実行する
 */
export async function flushAll(): Promise<void> {
  await flushPromises();
  jest.runAllTimers();
  await flushPromises();
}

/**
 * 指定されたイベントの発生回数を取得する
 * @param wrapper コンポーネントの Wrapper オブジェクト
 * @param eventName イベント名
 * @returns イベント発生回数
 */
export function getEventCount<T extends Vue | null>(
  wrapper: Wrapper<T>,
  eventName: string
): number {
  const eventInfo = wrapper.emitted()[eventName];
  if (!eventInfo) {
    return 0;
  }
  return eventInfo.length;
}

/**
 * 指定されたイベントの最初に発生したイベントデータ(配列)のうち、
 * 1つ目を型パラメータU型の値として取得する
 * @param wrapper コンポーネントの Wrapper オブジェクト
 * @param eventName イベント名
 * @returns イベントデータ
 */
export function getVeryFirstEventData<T extends Vue | null, U>(
  wrapper: Wrapper<T>,
  eventName: string
): U {
  const eventInfo = wrapper.emitted()[eventName];
  if (!eventInfo) {
    throw new Error(`イベントデータが存在しません(${eventName})`);
  }
  return eventInfo[0][0] as U;
}
