<template>
  <v-app>
    <v-app-bar app color="primary" dark>
      <div class="d-flex align-center">
        <v-img
          alt="Vuetify Logo"
          class="shrink mr-2"
          contain
          src="https://cdn.vuetifyjs.com/images/logos/vuetify-logo-dark.png"
          transition="scale-transition"
          width="40"
        />

        <v-img
          alt="Vuetify Name"
          class="shrink mt-1 hidden-sm-and-down"
          contain
          min-width="100"
          src="https://cdn.vuetifyjs.com/images/logos/vuetify-name-dark.png"
          width="100"
        />
      </div>

      <v-spacer></v-spacer>

      <v-btn href="https://github.com/vuetifyjs/vuetify/releases/latest" target="_blank" text>
        <span class="mr-2">Latest Release</span>
        <v-icon>mdi-open-in-new</v-icon>
      </v-btn>
    </v-app-bar>

    <v-main>
      <v-container>
        <!-- メッセージ -->
        <v-row v-if="message" justify="center">
          <v-col :md="msgColSize">
            <app-message :message-type="messageType" :message="message" />
          </v-col>
        </v-row>

        <router-view @msg="showMessage" />
      </v-container>
    </v-main>
  </v-app>
</template>

<script lang="ts">
import { Vue, Component, Watch } from "vue-property-decorator";
import AppMessage from "@/components/AppMessage.vue";
import { Message, MessageType } from "@/utils/message";
import { AppMsgSize } from "@/utils/consts";

@Component({
  components: {
    AppMessage,
  },
})
export default class App extends Vue {
  private messageType!: MessageType;
  private message = "";
  private msgColSize = AppMsgSize.Col6;

  /**
   * 画面上部のメッセージ表示領域にメッセージを表示する。
   * @param message メッセージ
   * @param msgColSize メッセージ表示サイズ
   */
  private showMessage(message: Message, colSize?: number) {
    this.messageType = message.messageType;
    this.message = message.message;
    this.msgColSize = colSize || AppMsgSize.Col6;
  }

  /**
   * $route オブジェクトのウォッチャー。
   * ページ遷移時にメッセージをクリアする。
   */
  @Watch("$route")
  private onRouteChanged() {
    this.message = "";
  }
}
</script>
