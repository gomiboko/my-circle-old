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
                <v-btn ref="loginButton" @click="login" :loading="loading" :disabled="loading" color="primary" block
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
import { ValidationObserver, ValidationProvider, extend, localize } from "vee-validate";
import { required } from "vee-validate/dist/rules";
import ja from "vee-validate/dist/locale/ja.json";
import SmallLink from "@/components/SmallLink.vue";
import { AppMessageSize } from "@/utils/app-message";

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
  private loading = false;

  private created() {
    this.$appMsg.setSize(AppMessageSize.Small);
  }

  private async login() {
    this.loading = true;

    const observer = this.$refs.observer as InstanceType<typeof ValidationObserver>;
    observer.reset();

    if (!(await observer.validate())) {
      this.loading = false;
      return;
    }

    try {
      await this.$http.post(
        "/login",
        {
          email: this.email,
          password: this.password,
        },
        {
          withCredentials: true,
        }
      );
      this.$router.push("/");
    } finally {
      this.loading = false;
    }
  }

  private gotoJoin() {
    this.$router.push("/join");
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
