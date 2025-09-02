const path = require("path");
const HtmlWebpackPlugin = require("html-webpack-plugin");
const CopyWebpackPlugin = require("copy-webpack-plugin");

module.exports = {
  entry: "./frontend/app.js",
  output: {
    path: path.resolve(__dirname, "dist"),
    filename: "bundle.js",
    clean: true,
  },
  module: {
    rules: [
      {
        test: /\.css$/i,
        use: ["style-loader", "css-loader"],
      },
    ],
  },
  plugins: [
    new HtmlWebpackPlugin({
      template: "./frontend/index.html",
      filename: "index.html",
    }),
    new CopyWebpackPlugin({
      patterns: [
        {
          from: path.resolve(__dirname, "wasm_exec.js"),
          to: "wasm_exec.js",
          noErrorOnMissing: true,
        },
        {
          from: path.resolve(__dirname, "frontend/styles.css"),
          to: "styles.css",
          noErrorOnMissing: true,
        },
        {
          from: path.resolve(__dirname, "dist/wasm"),
          to: "wasm",
          noErrorOnMissing: true,
        },
      ],
    }),
  ],
  devServer: {
    static: {
      directory: path.join(__dirname, "dist"),
    },
    compress: true,
    port: 8080,
  },
};
