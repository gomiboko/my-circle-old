import { createLocalVue } from "@vue/test-utils";
import axios from "axios";
import { VueConstructor } from "vue";
import VueRouter from "vue-router";

// axiosモジュールのモック化
jest.mock("axios");

/**
 * モックが設定されたローカルVueオブジェクトを生成する
 * @returns ローカルVueオブジェクトとモックオブジェクト
 */
export function createMockedLocalVue(): {
  localVue: VueConstructor<Vue>;
  axiosMock: jest.Mocked<typeof axios>;
} {
  const localVue = createLocalVue();

  // VueRouterの設定
  localVue.use(VueRouter);

  // axiosのモック設定
  const axiosMock = axios as jest.Mocked<typeof axios>;
  localVue.prototype.$http = axiosMock;

  return { localVue, axiosMock };
}
