const { createProxyMiddleware } = require('http-proxy-middleware');

module.exports = function(app) {
  app.use(
    '/webrtc/*',
    createProxyMiddleware({
      target: 'http://localhost:8081',
      secure: false,
      logLevel: 'debug',
      changeOrigin: true,
    })
  );
};
