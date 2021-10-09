<template>
  <div>
    <validation-observer v-slot="{ invalid }">
      <form>
        <!-- エラーメッセージ -->
        <v-row v-if="err" justify="center">
          <v-col md="6">
            <v-alert type="error" border="left" text dense>
              {{ err }}
            </v-alert>
          </v-col>
        </v-row>

        <v-row justify-md="center">
          <v-col md="4">
            <div class="text-center text-h5">My Circle にログイン</div>
          </v-col>
        </v-row>

        <v-row justify="center">
          <v-col md="4" lg="3" xl="2">
            <v-card outlined class="pa-4">
              <v-row>
                <v-col>
                  <validation-provider
                    rules="required"
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
              <v-row>
                <v-col>
                  <validation-provider
                    rules="required"
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
              <v-row class="mt-0">
                <v-col class="pt-0">
                  <!-- TODO: リンク先 -->
                  <router-link to="/" class="link text-body-2"
                    >パスワードを忘れた場合</router-link
                  >
                </v-col>
              </v-row>
              <v-row>
                <v-col>
                  <v-btn :disabled="invalid" @click="login" block
                    >ログイン</v-btn
                  >
                </v-col>
              </v-row>
            </v-card>
          </v-col>
        </v-row>
      </form>
    </validation-observer>

    <v-row justify-md="center" class="mt-4">
      <v-col md="3" class="text-center">
        <!-- TODO: リンク先 -->
        <router-link to="/" class="link text-body-2"
          >新規アカウント登録</router-link
        >
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
import { required } from "vee-validate/dist/rules";
import ja from "vee-validate/dist/locale/ja.json";
import axios from "axios";

extend("required", required);
localize("ja", ja);

@Component({
  components: {
    ValidationObserver,
    ValidationProvider,
  },
})
export default class Login extends Vue {
  private err = "";
  private email = "";
  private password = "";

  private async login() {
    const baseUrl = process.env.VUE_APP_BACKEND_BASE_URL;
    try {
      await axios.post(
        `${baseUrl}/login`,
        {
          email: this.email,
          password: this.password,
        },
        {
          withCredentials: true,
        }
      );
      // TODO: トップページに遷移
      this.$router.push("/");
    } catch (e) {
      // TODO: エラーメッセージの表示方法
      if (axios.isAxiosError(e) && e.response && e.response.data) {
        this.err = e.response.data.message;
      } else {
        this.err = `予期せぬエラー(${e})`;
      }
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
