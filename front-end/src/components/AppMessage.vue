<template>
  <v-alert :type="messageType" border="left" text dense>
    {{ message }}
  </v-alert>
</template>

<script lang="ts">
import { Vue, Component, Prop } from "vue-property-decorator";
import { MessageType } from "@/utils/message";

@Component
export default class AppMessage extends Vue {
  @Prop() private readonly message!: string;
  @Prop() private readonly messageType!: MessageType;
}
</script>
