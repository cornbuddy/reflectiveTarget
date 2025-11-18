const { readdirSync } = require('fs');
const path = require("path");

const CssMinimizerPlugin = require("css-minimizer-webpack-plugin");
const MiniCssExtractPlugin = require("mini-css-extract-plugin");
const TerserPlugin = require('terser-webpack-plugin');

module.exports = {
    mode: "production",
    entry: [
        "./js/index.js",
        ...findFilesWithExtension("./css", ".css"),
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
        }],
    },
    optimization: {
        minimize: true,
        minimizer: [
            new TerserPlugin(),
            new CssMinimizerPlugin(),
        ],
    },
    plugins: [new MiniCssExtractPlugin()],
};

function findFilesWithExtension(dir, extension) {
    return readdirSync(dir)
        .filter(file => file.endsWith(extension))
        .map(file => `${dir}/${file}`);
}
