module.exports = {
  transpileDependencies: ["vuetify"],
  devServer: {
    client: {
      webSocketURL: "ws://0.0.0.0:80/ws",
    },
  },
};
