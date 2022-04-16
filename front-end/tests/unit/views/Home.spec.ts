import { shallowMount } from "@vue/test-utils";
import Home from "@/views/Home.vue";
import { createMockedLocalVue } from "../local-vue";
import { User } from "@/responses/user";
import { Circle } from "@/responses/circle";
import flushPromises from "flush-promises";
import { AxiosResponse } from "axios";
import { AppMessageType } from "@/utils/app-message";
import { errorHandler } from "@/utils/global-error-handler";
import { initAppMsg } from "../test-utils";
import { messages } from "../test-consts";

jest.useFakeTimers();

beforeEach(() => {
  initAppMsg();
});

describe("Home.vue", () => {
  describe("初期表示", () => {
    test("ユーザ情報ロード中の場合", async () => {
      const { localVue, axiosMock } = createMockedLocalVue();

      const resCircles: Circle[] = [];
      const resUser = {
        Circles: resCircles,
      } as User;

      axiosMock.get.mockResolvedValue({
        data: { user: resUser },
      } as AxiosResponse);

      const wrapper = shallowMount(Home, { localVue });

      // 非同期処理完了前の状態を確認
      expect(wrapper.html()).toContain("loading...");

      await flushPromises();
    });

    test("サークルに所属していないユーザの場合", async () => {
      const { localVue, axiosMock } = createMockedLocalVue();

      const resCircles: Circle[] = [];
      const resUser = {
        Circles: resCircles,
      } as User;

      axiosMock.get.mockResolvedValue({
        data: { user: resUser },
      } as AxiosResponse);

      const wrapper = shallowMount(Home, { localVue });
      await flushPromises();

      expect(wrapper.html()).toContain("まだサークルに参加していません");
    });

    test("1つのサークルに所属しているユーザの場合", async () => {
      const { localVue, axiosMock } = createMockedLocalVue();

      const resCircles: Circle[] = [{ Name: "circle1" } as Circle];
      const resUser = {
        Circles: resCircles,
      } as User;

      axiosMock.get.mockResolvedValue({
        data: { user: resUser },
      } as AxiosResponse);

      const wrapper = shallowMount(Home, { localVue });
      await flushPromises();

      const circleCntWrapper = wrapper.findComponent({ ref: "circleCount" });

      expect(circleCntWrapper.text()).toBe("参加サークル数：1");
      expect(wrapper.html()).toContain("circle1");
    });

    test("複数のサークルに所属しているユーザの場合", async () => {
      const { localVue, axiosMock } = createMockedLocalVue();

      const resCircles: Circle[] = [{ Name: "circle1" } as Circle, { Name: "circle2" } as Circle];
      const resUser = {
        Circles: resCircles,
      } as User;

      axiosMock.get.mockResolvedValue({
        data: { user: resUser },
      } as AxiosResponse);

      const wrapper = shallowMount(Home, { localVue });
      await flushPromises();

      const circleCntWrapper = wrapper.findComponent({ ref: "circleCount" });

      expect(circleCntWrapper.text()).toBe("参加サークル数：2");
      expect(wrapper.html()).toContain("circle1");
      expect(wrapper.html()).toContain("circle2");
    });

    test("予期せぬエラーが発生した場合", async () => {
      const { localVue, axiosMock } = createMockedLocalVue();
      localVue.config.errorHandler = errorHandler;

      axiosMock.get.mockRejectedValue(new Error("エラーテスト"));
      axiosMock.isAxiosError.mockReturnValue(false);

      const wrapper = shallowMount(Home, { localVue });
      await flushPromises();

      expect(wrapper.vm.$data["loading"]).toBe(false);
      expect(wrapper.vm.$appMsg.type).toBe(AppMessageType.Error);
      expect(wrapper.vm.$appMsg.message).toBe(messages.UnexpectedErrorHasOccurred);
      expect(wrapper.html()).toContain("failed to load.");
    });
  });
});
