<template>
  <div>
    <form>
      <!-- Justify -->
      <v-row justify-md="center">
        <v-col md="3">
          <div class="text-center">My Circle にログイン</div>
        </v-col>
      </v-row>
      <v-row justify-md="center">
        <v-col md="3">
          <validation-provider rules="required|email|max:254" name="メールアドレス" v-slot="{ errors }">
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
          <validation-provider rules="required|alpha_num|min:8|max:64" name="パスワード" v-slot="{ errors }">
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
          <router-link to="/" class="link">パスワードを忘れた場合</router-link>
        </v-col>
      </v-row>
      <v-row justify-md="center">
        <v-col md="3">
          <v-btn block>ログイン</v-btn>
        </v-col>
      </v-row>
    </form>

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
import { ValidationProvider, extend, localize } from "vee-validate";
import { alpha_num, email, max, min, required } from "vee-validate/dist/rules";
import ja from "vee-validate/dist/locale/ja.json";

extend('required', required)
extend('alpha_num', alpha_num)
extend('email', email)
extend('min', min)
extend('max', max)
localize('ja', ja)

@Component({
  components: {
    ValidationProvider
  }
})
export default class Login extends Vue {
  private email = ''
  private password = ''
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
