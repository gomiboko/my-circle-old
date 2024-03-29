<template>
  <div>
    <validation-observer ref="observer">
      <form @keypress.enter="login">
        <v-row justify="center">
          <v-col md="4" lg="3" xl="2">
            <div class="text-center text-h5">My Circle にログイン</div>
          </v-col>
        </v-row>

        <v-row justify="center">
          <v-col md="4" lg="3" xl="2">
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
                <v-btn
                  ref="loginButton"
                  @click="login"
                  :loading="$state.loading"
                  :disabled="$state.loading"
                  color="primary"
                  block
                  >ログイン</v-btn
                >
              </v-col>
            </v-row>
          </v-col>
        </v-row>
      </form>
    </validation-observer>

    <v-row justify-md="center" class="mt-4">
      <v-col md="4" lg="3" xl="2">
        <div class="text-divider text-body-2">または</div>
      </v-col>
    </v-row>

    <v-row justify-md="center" class="mt-4">
      <v-col md="4" lg="3" xl="2">
        <v-btn ref="registerAccountButton" @click="gotoJoin" color="primary" outlined block>新規アカウント登録</v-btn>
      </v-col>
    </v-row>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import { ValidationObserver, ValidationProvider } from "vee-validate";
import SmallLink from "@/components/SmallLink.vue";
import { AppMessageSize } from "@/store/app-message";
import { API_PATHS, PAGE_PATHS } from "@/utils/consts";

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

  private created() {
    this.$state.appMsg.setSize(AppMessageSize.Small);
  }

  private async login() {
    const observer = this.$refs.observer as InstanceType<typeof ValidationObserver>;
    observer.reset();

    if (!(await observer.validate())) {
      return;
    }

    await this.$http.post(
      API_PATHS.SESSIONS,
      {
        email: this.email,
        password: this.password,
      },
      {
        withCredentials: true,
      }
    );
    this.$router.push(PAGE_PATHS.HOME);
  }

  private gotoJoin() {
    this.$router.push(PAGE_PATHS.JOIN);
  }
}
</script>

<style scoped>
.text-divider {
  --text-divider-gap: 0.2rem;
  display: flex;
  align-items: center;
  color: gray;
}
.text-divider::before,
.text-divider::after {
  content: "";
  height: 1px;
  background-color: silver;
  flex-grow: 1;
}
.text-divider::before {
  margin-right: var(--text-divider-gap);
}
.text-divider::after {
  margin-left: var(--text-divider-gap);
}
</style>
