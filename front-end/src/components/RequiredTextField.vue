<template>
  <v-text-field v-on="$listeners" v-bind="$attrs">
    <template v-slot:label
      >{{ label }} <sup><span style="color: crimson">*</span></sup></template
    >
  </v-text-field>
</template>

<script lang="ts">
import { Vue, Component, Prop } from "vue-property-decorator";

@Component({ inheritAttrs: false })
export default class RequiredTextField extends Vue {
  @Prop() private readonly label!: string;
}
</script>
