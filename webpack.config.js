const { join } = require('path')

module.exports = {
  entry: './public/js/src/main.js',
  output: {
    path: join(__dirname, '/public/js/dist/'),
    filename: 'bundle.js'
  },
  mode: 'development',
  performance: {
    hints: false,
  },
  module: {
    rules: [
      {
        test: /\.js$/,
        loader: 'babel-loader',
        exclude: /node_modules/,
        query: {
          'presets': ['env', 'stage-0']
        }
      }
    ]
  }
}
