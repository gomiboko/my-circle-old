<template>
  <div>
    <validation-observer v-slot="{ invalid }">
      <form>
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
                    ref="emailTextFieldProvider"
                    rules="required"
                    name="メールアドレス"
                    v-slot="{ errors }"
                  >
                    <v-text-field
                      ref="emailTextField"
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
                    ref="passwordTextFieldProvider"
                    rules="required"
                    name="パスワード"
                    v-slot="{ errors }"
                  >
                    <v-text-field
                      ref="passwordTextField"
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
                  <small-link text="パスワードを忘れた場合" to="/" />
                </v-col>
              </v-row>
              <v-row>
                <v-col>
                  <v-btn ref="loginButton" :disabled="invalid" @click="login" block>ログイン</v-btn>
                </v-col>
              </v-row>
            </v-card>
          </v-col>
        </v-row>
      </form>
    </validation-observer>

    <v-row justify-md="center" class="mt-4">
      <v-col md="3" class="text-center">
        <small-link text="新規アカウント登録" to="/join" />
      </v-col>
    </v-row>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import { ValidationObserver, ValidationProvider, extend, localize } from "vee-validate";
import { required } from "vee-validate/dist/rules";
import ja from "vee-validate/dist/locale/ja.json";
import SmallLink from "@/components/SmallLink.vue";
import { showError } from "@/utils/message";
import { AppMsgSize } from "@/utils/consts";

extend("required", required);
localize("ja", ja);

@Component({
  components: {
    ValidationObserver,
    ValidationProvider,
    SmallLink,
  },
})
export default class Login extends Vue {
  private email = "";
  private password = "";

  private async login() {
    const baseUrl = process.env.VUE_APP_BACKEND_BASE_URL;
    try {
      await this.$http.post(
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
      showError(this, e, AppMsgSize.Col4);
    }
  }
}
</script>
