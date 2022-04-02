import { shallowMount } from "@vue/test-utils";
import Home from "@/views/Home.vue";
import { createMockedLocalVue } from "../local-vue";
import { User } from "@/responses/user";
import { Circle } from "@/responses/circle";
import flushPromises from "flush-promises";
import { AxiosResponse } from "axios";
import { getEventCount, getVeryFirstEventData } from "../test-utils";
import { Message, MSG_EVENT } from "@/utils/message";

jest.useFakeTimers();

describe("Home.vue", () => {
  describe("初期表示", () => {
    test("ユーザ情報ロード中の場合", async () => {
      const { localVue } = createMockedLocalVue();
      const wrapper = shallowMount(Home, { localVue });

      expect(wrapper.html()).toContain("loading...");
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

    test("バックエンドとの通信に失敗した場合", async () => {
      const { localVue, axiosMock } = createMockedLocalVue();

      axiosMock.get.mockRejectedValue(new Error("エラーテスト"));
      axiosMock.isAxiosError.mockReturnValue(false);

      const wrapper = shallowMount(Home, { localVue });
      await flushPromises();

      expect(getEventCount(wrapper, MSG_EVENT)).toBe(1);
      const eventData = getVeryFirstEventData<Home, Message>(wrapper, MSG_EVENT);
      expect(eventData.message).toContain("予期せぬエラー");
      expect(wrapper.html()).toContain("failed to load.");
    });
  });
});
