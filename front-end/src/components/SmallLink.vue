<template>
  <router-link :to="to" class="link text-body-2">{{ text }}</router-link>
</template>

<script lang="ts">
import { Vue, Component, Prop } from "vue-property-decorator";

@Component
export default class SmalLink extends Vue {
  @Prop() private readonly to!: string;
  @Prop() private readonly text!: string;
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
