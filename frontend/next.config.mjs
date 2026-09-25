/** @type {import('next').NextConfig} */

function apiProxyTarget() {
  if (process.env.API_PROXY_TARGET === undefined) {
    return 'http://localhost:8080';
  }

  return process.env.API_PROXY_TARGET.trim();
}

const nextConfig = {
  allowedDevOrigins: ['10.98.0.41'],
  images: {
    remotePatterns: [
      {
        protocol: 'http',
        hostname: 'localhost',
        port: '8080',
        pathname: '/uploads/**',
      },
    ],
  },

  async redirects() {
    return [
      {
        source: '/',
        destination: '/posts',
        permanent: true,
      },
    ];
  },

  rewrites() {
    const target = apiProxyTarget();

    if (!target) {
      return [];
    }

    return [
      {
        source: '/api/:path*',
        destination: `${target}/:path*`,
      },
    ];
  },
};

export default nextConfig;
