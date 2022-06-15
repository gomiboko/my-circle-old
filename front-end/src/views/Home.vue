<template>
  <div>
    <div v-if="loading">
      <!-- TODO: スピナー、スケルトンローダー等表示する -->
      loading...
    </div>
    <!-- 通信に失敗した場合 -->
    <div v-else-if="me === null">
      <!-- TODO: 失敗したとき用のアイコンか何か表示する -->
      failed to load.
    </div>
    <!-- サークル未参加の場合 -->
    <div v-else-if="me.Circles.length === 0">
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
      <div ref="circleCount">{{ "参加サークル数：" + me.Circles.length }}</div>
      <li v-for="c in me.Circles" :key="c.ID">
        {{ c.Name }}
      </li>
    </div>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import { User } from "@/responses/user";
import { AppMessageSize } from "@/store/app-message";
import { API_PATHS, PAGE_PATHS } from "@/utils/consts";

@Component
export default class Home extends Vue {
  // ストアのロード中フラグだとビューの切り替えが適切にできない為、個別に管理する
  // (axiosの通信完了時にフラグが下ろされるが、その時点ではmeにユーザ情報が設定されていない)
  private loading = true;
  private me: User | null = null;

  private async created() {
    this.$state.appMsg.setSize(AppMessageSize.Medium);

    try {
      const res = await this.$http.get(API_PATHS.USERS_ME, { withCredentials: true });
      this.me = res.data.user as User;
    } finally {
      this.loading = false;
    }
  }

  private gotoCircleRegister() {
    this.$router.push(PAGE_PATHS.CIRCLE_REGISTER);
  }
}
</script>
