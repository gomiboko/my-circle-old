<template>
  <div>
    <validation-observer ref="observer">
      <form @keypress.enter.prevent="register">
        <v-row justify="center">
          <v-col md="4" lg="3" xl="2">
            <div class="text-center text-h5">サークル作成</div>
          </v-col>
        </v-row>

        <v-row justify="center">
          <v-col md="4" lg="3" xl="2">
            <v-row>
              <v-col>
                <validation-provider name="サークル名" rules="required|max:45" v-slot="{ errors }">
                  <required-text-field label="サークル名" v-model="circleName" :error-messages="errors">
                    <template v-slot:prepend>
                      <v-icon v-if="circleIconFile === null" size="32" @click="showFileChooser">
                        mdi-account-supervisor-circle
                      </v-icon>
                      <v-badge
                        v-else
                        @click.native="removeIcon"
                        icon="mdi-close-thick"
                        color="red lighten-2"
                        overlap
                        bordered
                        style="cursor: pointer"
                      >
                        <v-avatar size="32" style="cursor: pointer" @click.stop="showFileChooser">
                          <img :src="circleIconDataUrl" />
                        </v-avatar>
                      </v-badge>
                    </template>
                  </required-text-field>
                </validation-provider>
                <input type="file" ref="fileInput" @change="changeIcon" accept=".png,.jpg" style="display: none" />
              </v-col>
            </v-row>
            <v-row>
              <v-col class="text-right">
                <v-btn @click="cancel" outlined style="width: 100px" class="mr-2"> キャンセル </v-btn>
                <v-btn
                  @click="register"
                  :loading="$state.loading"
                  :disabled="$state.loading"
                  color="primary"
                  depressed
                  style="width: 100px"
                >
                  作成
                </v-btn>
              </v-col>
            </v-row>
          </v-col>
        </v-row>
      </form>
    </validation-observer>
  </div>
</template>

<script lang="ts">
import { AppMessageSize, AppMessageType } from "@/store/app-message";
import { Component, Vue } from "vue-property-decorator";
import { ValidationObserver, ValidationProvider } from "vee-validate";
import RequiredTextField from "@/components/RequiredTextField.vue";
import { API_PATHS, MAX_ICON_FILE_SIZE, MESSAGES, PAGE_PATHS } from "@/utils/consts";
import { Route, NavigationGuardNext } from "vue-router";
import { CONTENT_TYPE, createFormData } from "@/utils/http";
import { IsNotAllowedIconFileFormat } from "@/utils/image";

Component.registerHooks(["beforeRouteEnter"]);

@Component({
  components: {
    ValidationObserver,
    ValidationProvider,
    RequiredTextField,
  },
})
export default class CircleRegister extends Vue {
  private circleName = "";
  private circleIconFile: File | null = null;
  private circleIconDataUrl: string | null = null;
  private fromPath = "";

  private beforeRouteEnter(to: Route, from: Route, next: NavigationGuardNext) {
    next((vm) => {
      vm.$data["fromPath"] = from.path;
    });
  }

  private created() {
    this.$state.appMsg.setSize(AppMessageSize.Small);
  }

  private async register() {
    const observer = this.$refs.observer as InstanceType<typeof ValidationObserver>;
    observer.reset();

    if (!(await observer.validate())) {
      return;
    }

    const data = createFormData({
      circleName: this.circleName,
      circleIconFile: this.circleIconFile,
    });

    const res = await this.$http.post(API_PATHS.CIRCLES, data, {
      withCredentials: true,
      headers: CONTENT_TYPE.MULTIPART_FORM_DATA,
    });

    if (res.data.message) {
      this.$state.nextScreenMsg.type = res.data.messageType;
      this.$state.nextScreenMsg.message = res.data.message;
    }

    // FIXME: 作成したサークルの詳細ページに遷移。そこからメンバー招待等。
    this.$router.push(PAGE_PATHS.HOME);
  }

  private cancel() {
    this.$router.push(this.fromPath);
  }

  private showFileChooser() {
    const fileInput = this.$refs.fileInput as HTMLInputElement;
    fileInput.click();
  }

  private changeIcon() {
    const fileInput = this.$refs.fileInput as HTMLInputElement;
    const fileList = fileInput.files as FileList;
    const selectedFile = fileList.item(0);

    // 画像選択がキャンセルされた場合、何もしない
    // (Input要素のようにファイル未選択にはしない)
    if (selectedFile === null) {
      return;
    }

    this.$state.appMsg.message = "";

    // アイコンファイルのサイズチェック
    if (selectedFile.size > MAX_ICON_FILE_SIZE) {
      this.$state.appMsg.type = AppMessageType.Error;
      this.$state.appMsg.message = MESSAGES.OVER_MAX_ICON_FILE_SIZE;
      return;
    }

    // TODO: 画像の範囲選択ダイアログを表示

    // 選択されたファイルの読み込み
    const reader = new FileReader();
    reader.onload = (e) => {
      // FileReader.readAsDataURLで読み込んだ場合、resultは文字列
      const dataUrl = e.target?.result as string | null | undefined;

      if (!dataUrl) {
        this.$state.appMsg.type = AppMessageType.Error;
        this.$state.appMsg.message = MESSAGES.FAILED_TO_LOAD_FILE;
        return;
      }

      // アイコンファイルのフォーマットチェック
      if (IsNotAllowedIconFileFormat(dataUrl)) {
        this.$state.appMsg.type = AppMessageType.Error;
        this.$state.appMsg.message = MESSAGES.NOT_ALLOWED_ICON_FILE_FORMAT;
      } else {
        // サークルアイコンの表示を更新
        this.circleIconDataUrl = dataUrl;
        this.circleIconFile = selectedFile;
      }
    };
    reader.readAsDataURL(selectedFile);
  }

  private removeIcon() {
    this.circleIconDataUrl = null;
    this.circleIconFile = null;
  }
}
</script>
