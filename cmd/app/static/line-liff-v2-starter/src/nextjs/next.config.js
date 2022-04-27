module.exports = {
  assetPrefix: './',
  reactStrictMode: true,
  env: {
    LIFF_ID: process.env.LIFF_ID,
    LINE_NOTIFY_CLIENT_ID: process.env.LINE_NOTIFY_CLIENT_ID,
    LINE_NOTIFY_REDIRECT_URI: process.env.LINE_NOTIFY_REDIRECT_URI
  },
};
