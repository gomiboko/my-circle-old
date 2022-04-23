import { State } from "@/store/store";
import { AxiosStatic } from "axios";

declare module "vue/types/vue" {
  interface Vue {
    $http: AxiosStatic;
    $state: State;
  }
}
