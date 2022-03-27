<template>
  <div>
    <div v-if="loading">
      <!-- TODO: スピナー、スケルトンローダー等表示する -->
      loading...
    </div>
    <!-- 通信に失敗した場合 -->
    <div v-else-if="me === undefined">
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
          <v-btn color="primary" block>サークル作成</v-btn>
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
      <div>{{ "参加サークル数：" + me.Circles.length }}</div>
      <li v-for="c in me.Circles" :key="c.ID">
        {{ c.Name }}
      </li>
    </div>
  </div>
</template>

<script lang="ts">
import { showError } from "@/utils/message";
import { Component, Vue } from "vue-property-decorator";
import { User } from "@/responses/user";

@Component
export default class Home extends Vue {
  private loading = true;
  private me?: User;

  private async created() {
    try {
      const res = await this.$http.get(`${process.env.VUE_APP_BACKEND_BASE_URL}/users/me`, { withCredentials: true });
      this.me = res.data.user as User;
      this.loading = false;
    } catch (e) {
      showError(this, e);
      this.loading = false;
    }
  }
}
</script>
