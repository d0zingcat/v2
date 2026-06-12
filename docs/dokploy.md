# Dokploy Deployment

Use Dokploy's Dockerfile deployment mode instead of Docker Compose.

## Build Settings

- Repository: `d0zingcat/v2`
- Branch: `personal`
- Docker Context: `.`
- Dockerfile Path: `packaging/docker/alpine/Dockerfile`
- Internal Port: `8080`

## Environment Variables

```env
DATABASE_URL=user=miniflux password=your-password host=your-postgres-host dbname=miniflux sslmode=disable
RUN_MIGRATIONS=1
CREATE_ADMIN=1
ADMIN_USERNAME=d0zingcat
ADMIN_PASSWORD=your-admin-password
BASE_URL=https://miniflux.d0zingcat.dev
```

After the first successful deployment, `CREATE_ADMIN` can stay enabled because
Miniflux skips creation when the admin user already exists.
