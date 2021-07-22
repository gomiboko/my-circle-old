<template>
  <div>
    <v-form>
      <!-- Justify -->
      <v-row justify-md="center">
        <v-col md="3">
          <div class="text-center">My Circle にログイン</div>
        </v-col>
      </v-row>
      <v-row justify-md="center">
        <v-col md="3">
          <v-text-field
            label="メールアドレス"
          ></v-text-field>
        </v-col>
      </v-row>
      <v-row justify-md="center">
        <v-col md="3">
          <v-text-field
            label="パスワード"
            type="password"
          ></v-text-field>
        </v-col>
      </v-row>
      <v-row justify-md="center">
        <v-col md="3">
          <!-- TODO: リンク先 -->
          <router-link to="/" class="link">パスワードを忘れた場合</router-link>
        </v-col>
      </v-row>
      <v-row justify-md="center">
        <v-col md="3">
          <v-btn block>ログイン</v-btn>
        </v-col>
      </v-row>
    </v-form>

    <v-row>
      <v-col offset-md="5" md="2" class="text-center">
        <!-- TODO: リンク先 -->
        <router-link to="/" class="link">新規アカウント登録</router-link>
      </v-col>
    </v-row>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";

@Component
export default class Login extends Vue {
}
</script>

<style scoped>
.link {
  text-decoration: none;
}
.link:hover {
  text-decoration: underline;
}
</style>
