import Navbar from './components/Navbar';
import { AuthProvider } from './context/AuthContext';
import './globals.css';
import { ToastProvider } from './context/ToastContext';

export default function RootLayout({ children }) {
  return (
    <html lang="en">
      <body>
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