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

      <!-- FIXME: ログイン中のみ表示させる -->
      <v-menu offset-y nudge-bottom="5">
        <template v-slot:activator="{ on }">
          <v-avatar color="white" size="36" v-on="on" style="cursor: pointer" tabindex="-1">
            <!-- TODO: プロフィール画像が設定済みかどうかで分岐させる -->
            <v-icon light>mdi-account</v-icon>
          </v-avatar>
        </template>
        <v-list dense>
          <template v-for="(menu, index) in accountMenus">
            <v-list-item v-if="menu.NAME" @click="onMenuClick(menu.ID)" :key="index">
              <v-list-item-title>{{ menu.NAME }}</v-list-item-title>
            </v-list-item>
            <v-divider v-else :key="index"></v-divider>
          </template>
        </v-list>
      </v-menu>
    </v-app-bar>

    <v-main>
      <v-container>
        <!-- メッセージ -->
        <v-row v-if="$state.appMsg.message" justify="center">
          <v-col ref="appMessageColumn" :md="$state.appMsg.md" :lg="$state.appMsg.lg" :xl="$state.appMsg.xl">
            <app-message :messageType="$state.appMsg.type" :message="$state.appMsg.message" />
          </v-col>
        </v-row>

        <router-view />
      </v-container>
    </v-main>
  </v-app>
</template>

<script lang="ts">
import { Vue, Component } from "vue-property-decorator";
import AppMessage from "@/components/AppMessage.vue";

const ACCOUNT_MENU_ITEMS = Object.freeze({
  DIVIDER: Object.freeze({ ID: 0, NAME: "" }),
  PROFILE: Object.freeze({ ID: 1, NAME: "プロフィール" }),
  CONFIG: Object.freeze({ ID: 2, NAME: "設定" }),
  LOGOUT: Object.freeze({ ID: 3, NAME: "ログアウト" }),
});

@Component({
  components: {
    AppMessage,
  },
})
export default class App extends Vue {
  get accountMenus(): Readonly<{ ID: number; NAME: string }>[] {
    return [
      ACCOUNT_MENU_ITEMS.PROFILE,
      ACCOUNT_MENU_ITEMS.CONFIG,
      ACCOUNT_MENU_ITEMS.DIVIDER,
      ACCOUNT_MENU_ITEMS.LOGOUT,
    ];
  }

  private onMenuClick(id: number) {
    switch (id) {
      case ACCOUNT_MENU_ITEMS.PROFILE.ID:
        // TODO: プロフィール表示
        break;
      case ACCOUNT_MENU_ITEMS.CONFIG.ID:
        // TODO: 設定表示
        break;
      case ACCOUNT_MENU_ITEMS.LOGOUT.ID:
        // TODO: ログアウト処理
        break;
    }
  }
}
</script>
