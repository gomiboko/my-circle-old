<template>
  <div>
    <div v-if="loading">
      <!-- TODO: スピナー、スケルトンローダー等表示する -->
      loading...
    </div>
    <!-- 通信に失敗した場合 -->
    <div v-else-if="!userName">
      <!-- TODO: 失敗したとき用のアイコンか何か表示する -->
      failed to load.
    </div>
    <!-- サークル未参加の場合 -->
    <div v-else-if="circles.length === 0">
      <v-row justify="center" align="end" style="height: 150px">
        <v-col>
          <div class="text-center text-h6">まだサークルに参加していません</div>
        </v-col>
      </v-row>
      <v-row align="end" style="height: 100px">
        <v-col>
          <div class="text-center">サークルを作ろう！</div>
        </v-col>
      </v-row>
      <v-row justify="center">
        <v-col md="4" lg="3" xl="2">
          <v-btn @click="gotoCircleRegister" color="primary" block>サークル作成</v-btn>
        </v-col>
      </v-row>
      <v-row align="end" style="height: 100px">
        <v-col>
          <div class="text-center">サークルに参加しよう！</div>
        </v-col>
      </v-row>
      <v-row justify="center">
        <v-col md="4" lg="3" xl="2">
          <v-btn color="secondary" block>サークル参加</v-btn>
        </v-col>
      </v-row>
    </div>
    <!-- サークル情報が取得できた場合 -->
    <div v-else>
      <!-- TODO: 仮 -->
      <div ref="circleCount">{{ "参加サークル数：" + circles.length }}</div>
      <li v-for="c in circles" :key="c.ID">
        <v-avatar v-if="c.IconUrl" size="24">
          <img :src="c.IconUrl" />
        </v-avatar>
        <v-avatar v-else size="24" :style="createBgColorStyleFromText(c.Name)">
          <!-- サークル名の先頭1文字をアイコンとして使用する -->
          <span class="white--text text-h6">{{ c.Name.substr(0, 1) }}</span>
        </v-avatar>
        {{ c.Name }}
      </li>
    </div>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import { Circle } from "@/responses/circle";
import { AppMessageSize } from "@/store/app-message";
import { API_PATHS, PAGE_PATHS } from "@/utils/consts";
import { createBgColorStyleFromText } from "@/utils/view";

@Component
export default class Home extends Vue {
  // ストアのロード中フラグだとビューの切り替えが適切にできない為、個別に管理する
  // (axiosの通信完了時にフラグが下ろされるが、その時点ではuserNameにユーザ名が設定されていない)
  private loading = true;

  private userName = "";
  private userIconUrl = "";
  private circles: Circle[] = [];

  private async created() {
    this.$state.appMsg.setSize(AppMessageSize.Medium);

    try {
      const res = await this.$http.get(API_PATHS.USERS_ME, { withCredentials: true });

      this.userName = res.data.userName as string;
      this.userIconUrl = res.data.userIconUrl as string;
      if (res.data.circles) {
        this.circles.push(...res.data.circles);
      }
    } finally {
      this.loading = false;
    }
  }

  private gotoCircleRegister() {
    this.$router.push(PAGE_PATHS.CIRCLE_REGISTER);
  }

  private createBgColorStyleFromText(text: string): string {
    return createBgColorStyleFromText(text);
  }
}
</script>
