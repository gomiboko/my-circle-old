<template>
  <v-alert :type="messageType" border="left" text dense>
    {{ message }}
  </v-alert>
</template>

<script lang="ts">
import { AppMessageType } from "@/store/app-message";
import { Vue, Component, Prop } from "vue-property-decorator";

@Component
export default class AppMessage extends Vue {
  @Prop() private readonly message!: string;
  @Prop() private readonly messageType!: AppMessageType;
}
</script>
