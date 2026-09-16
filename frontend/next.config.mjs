/** @type {import('next').NextConfig} */

const nextConfig = {
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
    if (process.env.NODE_ENV === 'development') {
      return [
        {
          source: '/api/:path*',
          destination: 'http://localhost:8080/:path*',
        },
      ];
    }

    return [];
  },
};  

export default nextConfig;
