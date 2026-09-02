import './globals.css';

export const metadata = {
  title: 'Social Network',
  description: 'My awesome social network project',
};

export default function RootLayout({ children }) {
  return (
    <html lang="en">
      <body>
        {children}
      </body>
    </html>
  );
}