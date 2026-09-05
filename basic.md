# Basic React concepts 
--------------------------------------
1. Components 
2. JSX 
3. Props 
4. State 
5. Events
6. Conditional Rendering
7. Rendering Lists 
8. useEffect 
9. Hooks 
10. Client vs Server Components
11. Next.js routing/data fetching
---------------------------------- 
# React vs Next.js: Detailed Comparison
 
A feature-by-feature comparison of React and Next.js to help you decide which fits your project.
 
| Feature | React | Next.js |
|---|---|---|
| **Type** | UI Library | Full-Stack Framework |
| **Rendering Strategy** | Client-Side Rendering (CSR): Renders HTML in the browser via JavaScript. Initial load can show a blank screen momentarily. | Hybrid (Server & Client): Supports Server-Side Rendering (SSR), Static Site Generation (SSG), and Server Components for faster initial loads. |
| **SEO & Crawlability** | Harder for search engines to index out of the box because the initial page is empty HTML populated by JS. | Excellent for SEO because fully rendered HTML is served directly from the server to web crawlers. |
| **Routing** | Requires external libraries like `react-router-dom` to configure URL routes manually. | File-system-based routing: Simply creating a file or directory (e.g., `app/about/page.js`) creates the route automatically. |
| **Backend Integration** | Focuses purely on the frontend UI. Requires a separate backend server (e.g., Express, Go, Fastify) for API routes. | Supports built-in API routes and Server Actions, allowing you to build full-stack apps within a single codebase. |
| **Setup & Tooling** | Flexible, but requires picking your own build tool (e.g., Vite), state management, and router. | Opinionated and pre-configured out of the box with zero setup needed for bundling, routing, and compilation. |
 
## Summary
 
- **Choose React** if you want maximum flexibility and are building a client-heavy app (e.g., a dashboard) where you'll assemble your own toolchain.
- **Choose Next.js** if you need strong SEO, fast initial page loads, or want a full-stack solution without managing a separate backend.
 


---------------------------------
+ Server-Side Rendering (SSR) & Static Generation   |
|   + File-based Routing (App / Pages Router)           |
|   + Built-in API Routes & Server Actions              |
|   + Automatic Image / Font Optimizations