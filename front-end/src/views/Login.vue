<template>
  <div>
    <validation-observer v-slot="{ invalid }">
      <form>
        <v-row justify-md="center">
          <v-col md="3">
            <div class="text-center">My Circle にログイン</div>
          </v-col>
        </v-row>
        <v-row justify-md="center">
          <v-col md="3">
            <validation-provider
              rules="required|email|max:254"
              name="メールアドレス"
              v-slot="{ errors }"
            >
              <v-text-field
                label="メールアドレス"
                v-model="email"
                :error-messages="errors"
              ></v-text-field>
            </validation-provider>
          </v-col>
        </v-row>
        <v-row justify-md="center">
          <v-col md="3">
            <validation-provider
              rules="required|alpha_num|min:8|max:64"
              name="パスワード"
              v-slot="{ errors }"
            >
              <v-text-field
                label="パスワード"
                type="password"
                v-model="password"
                :error-messages="errors"
              ></v-text-field>
            </validation-provider>
          </v-col>
        </v-row>
        <v-row justify-md="center">
          <v-col md="3">
            <!-- TODO: リンク先 -->
            <router-link to="/" class="link"
              >パスワードを忘れた場合</router-link
            >
          </v-col>
        </v-row>
        <v-row justify-md="center">
          <v-col md="3">
            <v-btn :disabled="invalid" @click="login" block>ログイン</v-btn>
          </v-col>
        </v-row>
      </form>
    </validation-observer>

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
import {
  ValidationObserver,
  ValidationProvider,
  extend,
  localize,
} from "vee-validate";
import { alpha_num, email, max, min, required } from "vee-validate/dist/rules";
import ja from "vee-validate/dist/locale/ja.json";
import axios from "axios";

extend("required", required);
extend("alpha_num", alpha_num);
extend("email", email);
extend("min", min);
extend("max", max);
localize("ja", ja);

@Component({
  components: {
    ValidationObserver,
    ValidationProvider,
  },
})
export default class Login extends Vue {
  private email = "";
  private password = "";

  private async login() {
    const baseUrl = process.env.VUE_APP_BACKEND_BASE_URL;
    try {
      await axios.post(`${baseUrl}/login`, {
        email: this.email,
        password: this.password,
      });
      // TODO: トップページに遷移
      this.$router.push("/");
    } catch(e) {
      // TODO: エラーメッセージの表示方法
      alert("エラーが発生しました。");
    }
  }
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
