# social-network

Monorepo with three parts:

| Folder      | What it is                                          |
| ----------- | --------------------------------------------------- |
| `backend/`  | Go API + SQLite (listens on `:8080`, paths at root) |
| `frontend/` | Next.js web app (listens on `:3000`)                |
| `mobile/`   | Expo / React Native app                             |

## Running the API

```bash
docker compose up --build
```

Caddy publishes **`:8000`** and strips the `/api` prefix before proxying to the Go
server, so the API is reached at `http://<your-lan-ip>:8000/api/...`.

Check it is healthy:

```bash
curl http://localhost:8000/api/health      # -> {"status":"ok"}
```

Add the web UI (not needed for mobile development):

```bash
docker compose --profile web up --build
```

The Go server also accepts `/api/...` directly, so it works without Caddy:

```bash
cd backend && go run ./cmd/server          # serves on :8080
```

## Running the mobile app

The client targets `EXPO_PUBLIC_API_URL` (falling back to `10.0.2.2:8000` on
Android, `localhost:8000` elsewhere) and always appends `/api`.

### Android emulator

`10.0.2.2` is the emulator's alias for your host machine, so the default works
as long as the API is up on port 8000:

```bash
docker compose up --build
cd mobile && npx expo start --android
```

### Physical phone (Expo Go)

`10.0.2.2` does **not** exist on a real device — point the app at your computer's
LAN IP and restart Metro with `--clear` (the variable is inlined at bundle time):

```bash
hostname -I | awk '{print $1}'                       # e.g. 192.168.1.50
docker compose up --build

cd mobile
EXPO_PUBLIC_API_URL=http://192.168.1.50:8000 npx expo start --clear
```

The phone must be on the same network. If it still fails, open
`http://192.168.1.50:8000/api/health` in the phone's browser — if that does not
show `{"status":"ok"}`, the problem is the network/firewall, not the app.

Over USB, `adb reverse` avoids the LAN entirely:

```bash
adb reverse tcp:8000 tcp:8000
EXPO_PUBLIC_API_URL=http://localhost:8000 npx expo start --clear
```

> `npx expo start --tunnel` only tunnels the Metro bundler, not the API.
