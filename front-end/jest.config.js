module.exports = {
  preset: "@vue/cli-plugin-unit-jest/presets/typescript-and-babel",
  setupFiles: ["./tests/unit/setup.ts"],
  transform: {
    "vee-validate/dist/rules": "babel-jest",
  },
  transformIgnorePatterns: ["<rootDir>/node_modules/(?!vee-validate/dist/rules)"],
};
