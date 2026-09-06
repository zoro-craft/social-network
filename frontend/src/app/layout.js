import './globals.css';
import MouseGlow from '@/components/MouseGlow';
import { AuthProvider } from '@/app/context/AuthContext';
import { ToastProvider } from '@/app/context/ToastContext';
import Navbar from '@/components/Navbar';

export const metadata = {
  title: 'Social Network',
  description: 'My awesome social network project',
};

export default function RootLayout({ children }) {
  return (
    <html lang="en" suppressHydrationWarning>
      <body>
        <MouseGlow />
        <AuthProvider>
          <ToastProvider>
            <Navbar />
            {children}
          </ToastProvider>
        </AuthProvider>
      </body>
    </html>
  );
}