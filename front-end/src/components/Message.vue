<template>
  <v-alert :type="type" border="left" text dense>
    {{ message }}
  </v-alert>
</template>

<script lang="ts">
import { Vue, Component, Prop } from "vue-property-decorator";

@Component
export default class Message extends Vue {
  @Prop() private readonly message!: string;
  @Prop() private readonly type!: string;
}
</script>
