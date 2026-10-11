'use client';

import { usePathname } from 'next/navigation';
import LeftSidebar from './LeftSidebar';

// Routes that own the whole viewport instead of sitting in the content column:
// the auth screens are centered cards, and the chat page already brings its own
// two-column layout, so neither wants the shell's frame around it.
const FULL_BLEED = ['/login', '/register', '/chat'];

// Pages that are exactly one viewport tall and scroll internally. They still use
// the shell, but not its bottom padding, which would add a stray scroll.
const IMMERSIVE = ['/reels'];

// Desktop shell: a groups/online gutter and the page content.
export default function AppShell({ children }) {
  const pathname = usePathname();

  const fullBleed = FULL_BLEED.some(
    (route) => pathname === route || pathname.startsWith(`${route}/`),
  );

  const immersive = IMMERSIVE.some((route) => pathname.startsWith(route));

  if (fullBleed) return children;

  return (
    <div className={`glass-shell ${immersive ? 'glass-shell-immersive' : ''}`}>
      <LeftSidebar />
      {/* A div, not a main element: several pages render a <main> of their own and
          nesting the two would be invalid. */}
      <div className="glass-shell-main">{children}</div>
    </div>
  );
}
