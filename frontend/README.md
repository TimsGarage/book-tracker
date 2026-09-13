# Book Tracker Frontend (Tauri + SvelteKit)

Mobile and desktop frontend for Book Tracker built with [Tauri v2](https://v2.tauri.app/) and [SvelteKit](https://kit.svelte.dev/).

---

## Configuration & Backend URL

The frontend communicates with the Go backend API.

- **Local Development**: In dev mode (`npm run dev` or `npm run tauri dev`), Vite proxies `/api` requests directly to `http://127.0.0.1:8080`. No environment configuration is needed.
- **Android / Production**: The mobile APK runs directly on your device and connects across the network. The backend URL is baked into the build using the `VITE_API_URL` environment variable:
  ```bash
  # Example:
  VITE_API_URL="https://api.mybooktracker.com"
  ```
  See [.env.example](.env.example) for reference.

---

## Building Android APK via GitHub Actions Pipeline

A GitHub Actions workflow is provided at [`.github/workflows/frontend-android.yml`](../.github/workflows/frontend-android.yml) to automatically compile, sign, and release installable `.apk` files for personal use.

### Method 1: Manual Run (Interactive URL Input)
1. Go to your GitHub repository -> **Actions**.
2. Select **Build and Release Frontend Android APK**.
3. Click **Run workflow**:
   - **Backend API URL**: Enter your deployed backend address (e.g. `https://api.mybooktracker.com` or `http://192.168.1.50:8080`).
   - **Build Flavor**: Choose `release` (optimized) or `debug`.
   - **Target CPU Architectures**: Choose `all` (broad compatibility) or `aarch64` (modern 64-bit phones for faster build times).
4. When complete, download the APK from the workflow's **Artifacts** section and install it directly on your Android phone!

### Method 2: Automatic Release on Git Tag
1. Save your backend URL once as a GitHub Repository Variable:
   - Go to **Settings > Secrets and variables > Actions > Variables**.
   - Create a new variable named `VITE_API_URL` (e.g. `https://api.mybooktracker.com`).
2. Push a version tag:
   ```bash
   git tag v1.0.0
   git push origin v1.0.0
   ```
3. The workflow will automatically build the APKs, sign them, and create a GitHub Release with the `.apk` assets attached for instant download.

---

## Building Locally (Prerequisites)

- Node.js 20+
- Rust & Cargo
- Android Studio / Android SDK & NDK 26+
- Rust Android targets:
  ```bash
  rustup target add aarch64-linux-android armv7-linux-androideabi i686-linux-android x86_64-linux-android
  ```

Build APK with custom backend URL:
```bash
cd frontend
VITE_API_URL="https://api.mybooktracker.com" npm run tauri android build -- --apk
```
The output APK will be placed in `src-tauri/gen/android/app/build/outputs/apk/`.
