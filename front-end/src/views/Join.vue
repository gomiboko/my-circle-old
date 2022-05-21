<template>
  <div>
    <validation-observer ref="observer">
      <form @keypress.enter="register">
        <v-row justify="center">
          <v-col md="4" lg="3" xl="2">
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
                  name="パスワード"
                  rules="required|min:8|max:128|password"
                  v-slot="{ errors }"
                >
                  <required-text-field
                    ref="passwordTextField"
                    label="パスワード"
                    hint="8文字以上"
                    v-model="password"
                    :type="showPassword ? 'text' : 'password'"
                    :error-messages="errors"
                  >
                    <!-- v-text-field の append-icon プロパティだと tabindex が指定できないのでアイコンスロットを使用する -->
                    <template v-slot:append>
                      <v-icon @click="showPassword = !showPassword" tabindex="-1">
                        {{ showPassword ? "mdi-eye" : "mdi-eye-off" }}
                      </v-icon>
                    </template>
                  </required-text-field>
                </validation-provider>
              </v-col>
            </v-row>
            <v-row>
              <v-col>
                <v-btn
                  ref="registerButton"
                  @click="register"
                  :loading="$state.loading"
                  :disabled="$state.loading"
                  color="primary"
                  block
                  >登録する</v-btn
                >
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
import { ValidationObserver, ValidationProvider } from "vee-validate";
import RequiredTextField from "@/components/RequiredTextField.vue";
import { AppMessageSize } from "@/store/app-message";
import { API_PATHS, PAGE_PATHS } from "@/utils/consts";

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

  private created() {
    this.$state.appMsg.setSize(AppMessageSize.Small);
  }

  private async register() {
    const observer = this.$refs.observer as InstanceType<typeof ValidationObserver>;
    observer.reset();

    if (!(await observer.validate())) {
      return;
    }

    await this.$http.post(
      API_PATHS.USERS,
      {
        username: this.username,
        email: this.email,
        password: this.password,
      },
      {
        withCredentials: true,
      }
    );

    // トップページに遷移
    this.$router.push(PAGE_PATHS.HOME);
  }
}
</script>
