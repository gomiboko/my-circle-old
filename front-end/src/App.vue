<template>
  <v-app>
    <v-app-bar app color="primary" dark>
      <!-- TODO: アプリロゴ等に変更 -->
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

      <template v-if="isLoginRequiredPage">
        <v-menu ref="accountMenu" offset-y nudge-bottom="5">
          <template v-slot:activator="{ on }">
            <v-avatar color="white" size="36" v-on="on" style="cursor: pointer" tabindex="-1">
              <!-- TODO: プロフィール画像が設定済みかどうかで分岐させる -->
              <v-icon light>mdi-account</v-icon>
            </v-avatar>
          </template>
          <v-list dense>
            <template v-for="(menu, index) in accountMenus">
              <v-list-item v-if="menu.NAME" @click="onMenuClick(menu.ID)" :key="index">
                <v-list-item-title>
                  <v-icon color="blue-grey lighten-3" class="mr-2" dense>{{ menu.ICON }}</v-icon>
                  {{ menu.NAME }}
                </v-list-item-title>
              </v-list-item>
              <v-divider v-else :key="index"></v-divider>
            </template>
          </v-list>
        </v-menu>
      </template>
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
import { API_PATHS, PAGE_PATHS } from "./utils/consts";

const ACCOUNT_MENU_ITEMS = Object.freeze({
  DIVIDER: Object.freeze({ ID: 0, NAME: "", ICON: "" }),
  PROFILE: Object.freeze({ ID: 1, NAME: "プロフィール", ICON: "mdi-account-cog-outline" }),
  CONFIG: Object.freeze({ ID: 2, NAME: "設定", ICON: "mdi-cog-outline" }),
  LOGOUT: Object.freeze({ ID: 3, NAME: "ログアウト", ICON: "mdi-logout" }),
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

  get isLoginRequiredPage(): boolean {
    return ![PAGE_PATHS.LOGIN, PAGE_PATHS.JOIN].includes(this.$route.path);
  }

  private async onMenuClick(id: number) {
    switch (id) {
      case ACCOUNT_MENU_ITEMS.PROFILE.ID:
        // TODO: プロフィール表示
        break;
      case ACCOUNT_MENU_ITEMS.CONFIG.ID:
        // TODO: 設定表示
        break;
      case ACCOUNT_MENU_ITEMS.LOGOUT.ID:
        await this.$http.get(API_PATHS.SESSIONS, { withCredentials: true });
        this.$router.push(PAGE_PATHS.LOGIN);
        break;
    }
  }
}
</script>
