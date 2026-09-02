import { NextResponse } from 'next/server';

export function middleware(request) {
  const sessionToken = request.cookies.get('session_token');
  const { pathname } = request.nextUrl;

  if (sessionToken && (pathname === '/login' || pathname === '/register')) {
    return NextResponse.redirect(new URL('/posts', request.url));
  }

  return NextResponse.next();
}

export const config = {
  matcher: ['/login', '/register'],
};
