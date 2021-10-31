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

      <v-btn
        href="https://github.com/vuetifyjs/vuetify/releases/latest"
        target="_blank"
        text
      >
        <span class="mr-2">Latest Release</span>
        <v-icon>mdi-open-in-new</v-icon>
      </v-btn>
    </v-app-bar>

    <v-main>
      <v-container>
        <!-- メッセージ -->
        <v-row v-if="message" justify="center">
          <v-col md="6">
            <alert :message-type="messageType" :message="message" />
          </v-col>
        </v-row>

        <router-view @msg="showMessage" />
      </v-container>
    </v-main>
  </v-app>
</template>

<script lang="ts">
import { Vue, Component, Watch } from "vue-property-decorator";
import { Route } from "vue-router";
import Alert from "@/components/Alert.vue";
import { Message, MessageType } from "@/utils/message";

@Component({
  components: {
    Alert,
  },
})
export default class App extends Vue {
  private messageType!: MessageType;
  private message = "";

  /**
   * 画面上部のメッセージ表示領域にメッセージを表示する。
   */
  private showMessage(message: Message) {
    this.messageType = message.messageType;
    this.message = message.message;
  }

  /**
   * $route オブジェクトのウォッチャー。
   * ページ遷移時にメッセージをクリアする。
   */
  @Watch("$route")
  private onRouteChanged(to: Route, from: Route) {
    this.message = "";
  }
}
</script>
