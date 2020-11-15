const path = require('path');

const distPath = path.join(__dirname, 'dist');
const nmPath = path.join(__dirname, 'node_modules');
const entryPath = path.join(__dirname, 'src/js/index.js');
module.exports = {
    mode: 'development',
    entry: entryPath,
    output: {
        filename: 'bundle.js',
        path: distPath,
    },
    devServer: {
        contentBase: distPath,
        watchContentBase: true,
        compress: true,
        port: 9000,
        host: '0.0.0.0',
        index: 'index.html',
        writeToDisk: true
    },
    module: {
        rules: [
            {
                test: /\.html$/i,
                loader: 'html-loader',
            },
            {
                test: /\.m?js$/,
                exclude: nmPath,
                use: ['babel-loader', 'eslint-loader']
            }
        ],
    },
};
