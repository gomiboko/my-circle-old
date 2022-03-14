<template>
  <div>
    <validation-observer v-slot="{ invalid }">
      <form>
        <v-row justify-md="center">
          <v-col md="4">
            <div class="text-center text-h5">アカウント作成</div>
          </v-col>
        </v-row>

        <v-row justify="center">
          <v-col md="4" lg="3" xl="2">
            <v-row>
              <v-col>
                <validation-provider
                  ref="usernameTextFieldProvider"
                  name="ユーザ名"
                  rules="required|max:45"
                  v-slot="{ errors }"
                >
                  <required-text-field
                    ref="usernameTextField"
                    label="ユーザ名"
                    v-model="username"
                    :error-messages="errors"
                  ></required-text-field>
                </validation-provider>
              </v-col>
            </v-row>
            <v-row>
              <v-col>
                <validation-provider
                  ref="emailTextFieldProvider"
                  name="メールアドレス"
                  rules="required|email|max:254"
                  v-slot="{ errors }"
                >
                  <required-text-field
                    ref="emailTextField"
                    label="メールアドレス"
                    v-model="email"
                    :error-messages="errors"
                  ></required-text-field>
                </validation-provider>
              </v-col>
            </v-row>
            <v-row>
              <v-col>
                <validation-provider
                  ref="passwordTextFieldProvider"
                  rules="required|min:8|max:128|password"
                  name="パスワード"
                  v-slot="{ errors }"
                >
                  <required-text-field
                    ref="passwordTextField"
                    label="パスワード"
                    hint="8文字以上"
                    v-model="password"
                    :type="showPassword ? 'text' : 'password'"
                    :append-icon="showPassword ? 'mdi-eye' : 'mdi-eye-off'"
                    :error-messages="errors"
                    @click:append="showPassword = !showPassword"
                  ></required-text-field>
                </validation-provider>
              </v-col>
            </v-row>
            <v-row>
              <v-col>
                <v-btn ref="registerButton" :disabled="invalid" @click="register" color="primary" block>登録する</v-btn>
              </v-col>
            </v-row>
          </v-col>
        </v-row>
      </form>
    </validation-observer>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import { ValidationObserver, ValidationProvider, extend, localize } from "vee-validate";
import { required, min, max, email } from "vee-validate/dist/rules";
import ja from "vee-validate/dist/locale/ja.json";
import { showError } from "@/utils/message";
import RequiredTextField from "@/components/RequiredTextField.vue";
import { AppMsgSize } from "@/utils/consts";

extend("required", required);
extend("min", min);
extend("max", max);
extend("email", {
  ...email,
  // デフォルトのメッセージだと不自然になるので上書き
  message: "有効なメールアドレスではありません",
});
extend("password", {
  validate(value) {
    return /^[a-zA-Z0-9!@#$%^&*()-_=+[\]{}\\|~;:'",.<>/?`]*$/.test(value);
  },
  message: "{_field_}に使用できない文字が含まれています",
});
localize("ja", ja);

@Component({
  components: {
    ValidationObserver,
    ValidationProvider,
    RequiredTextField,
  },
})
export default class Join extends Vue {
  private username = "";
  private email = "";
  private password = "";
  private showPassword = false;

  private async register() {
    const baseUrl = process.env.VUE_APP_BACKEND_BASE_URL;
    try {
      await this.$http.post(
        `${baseUrl}/users`,
        {
          username: this.username,
          email: this.email,
          password: this.password,
        },
        {
          withCredentials: true,
        }
      );

      // トップページに繊維
      this.$router.push("/");
    } catch (e) {
      showError(this, e, AppMsgSize.Col4);
    }
  }
}
</script>
