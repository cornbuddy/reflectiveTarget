const CssMinimizerPlugin = require("css-minimizer-webpack-plugin");
const MiniCssExtractPlugin = require("mini-css-extract-plugin");
const path = require("path");

module.exports = {
    mode: "production",
    entry: [
        "./js/index.js",
        "./css/index.css",
    ],
    output: {
        filename: "[name].js",
        path: path.resolve(__dirname, "dist"),
    },
    resolve: {
        fallback: {
            "fs": false,
        },
    },
    module: {
        rules: [{
            test: /\.css$/i,
            use: [
                MiniCssExtractPlugin.loader,
                "css-loader",
            ],
        },],
    },
    optimization: {
        minimizer: [new CssMinimizerPlugin()],
    },
    plugins: [new MiniCssExtractPlugin()],
};
