const path = require('path')

module.exports = {
  entry: './public/js/src/main.js',
  output: {
    // eslint-disable-next-line no-undef
    path: path.join(__dirname, '/public/js/dist/'),
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
