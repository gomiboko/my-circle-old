import { AppMessage } from "@/utils/app-message";
import { AxiosStatic } from "axios";

declare module "vue/types/vue" {
  interface Vue {
    $http: AxiosStatic;
    $appMsg: AppMessage;
  }
}
